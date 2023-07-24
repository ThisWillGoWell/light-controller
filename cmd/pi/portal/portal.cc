#include <iostream>
#include <fstream>
#include <cstring>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <unistd.h>
#include "led-matrix.h"

constexpr int BUFFER_SIZE = 4096;

using rgb_matrix::RGBMatrix;
using rgb_matrix::Canvas;
using rgb_matrix::FrameCanvas;

RGBMatrix *matrix;
FrameCanvas *offscreen;

const int NUMBER_COLUMNS = 128;
const int NUMBER_ROWS = 96;
const int DATA_SIZE = NUMBER_COLUMNS * NUMBER_ROWS * 3;

void saveImage(const char* data, int dataSize)
{

    int pixelCounter = 0;
    // Copy all the pixels to the canvas.
    for (int y = 0; y < NUMBER_ROWS; ++y) {
        for (int x = 0; x < NUMBER_COLUMNS; ++x) {
          offscreen->SetPixel(x, y, uint8_t(data[pixelCounter*3]), uint8_t(data[pixelCounter*3+1]), uint8_t(data[pixelCounter*3+2]));
          pixelCounter++;
        }
    }
    offscreen = matrix->SwapOnVSync(offscreen);
}


void waiting_for_connection()
{
    offscreen->Fill(0, 0, 0);
    offscreen->SetPixel(0,NUMBER_ROWS-1, 50, 0, 0);
    offscreen = matrix->SwapOnVSync(offscreen);
}


void configureMatrix() {
    RGBMatrix::Options defaults;
    rgb_matrix::RuntimeOptions runtime_opt;
    runtime_opt.gpio_slowdown = 4;
    defaults.rows = 32;
    defaults.cols = 64;
    defaults.chain_length = 2;
    defaults.parallel = 3;
    defaults.show_refresh_rate = false;

    matrix = RGBMatrix::CreateFromOptions(defaults, runtime_opt);
    offscreen = matrix->CreateFrameCanvas();
    std::cout << "Created Matrix: " << matrix->width() << "x" << matrix->height() << std::endl;

}
int main()
{
    configureMatrix();
 int serverSocket = socket(AF_INET, SOCK_STREAM, 0);
    if (serverSocket == -1)
    {
        std::cerr << "Failed to create socket" << std::endl;
        return 1;
    }

    sockaddr_in serverAddress{};
    serverAddress.sin_family = AF_INET;
    serverAddress.sin_addr.s_addr = INADDR_ANY;
    serverAddress.sin_port = htons(8080);

    if (bind(serverSocket, reinterpret_cast<sockaddr*>(&serverAddress), sizeof(serverAddress)) == -1)
    {
        std::cerr << "Failed to bind socket to port" << std::endl;
        close(serverSocket);
        return 1;
    }

    if (listen(serverSocket, 1) == -1)
    {
        std::cerr << "Failed to listen on socket" << std::endl;
        close(serverSocket);
        return 1;
    }

    std::cout << "Waiting for a connection on port 8080..." << std::endl;

    while (true)
    {
        waiting_for_connection();
        sockaddr_in clientAddress{};
        socklen_t clientAddressSize = sizeof(clientAddress);

        int clientSocket = accept(serverSocket, reinterpret_cast<sockaddr*>(&clientAddress), &clientAddressSize);
        if (clientSocket == -1)
        {
            std::cerr << "Failed to accept client connection" << std::endl;
            continue;
        }

        std::cout << "Client connected. Receiving images..." << std::endl;
        char imageData[DATA_SIZE];
        int totalBytesReceived = 0;
        while (true) {
            totalBytesReceived = 0;
            while (totalBytesReceived < DATA_SIZE)
            {
                int bytesReceived = recv(clientSocket, imageData + totalBytesReceived, DATA_SIZE - totalBytesReceived, 0);
                if (bytesReceived <= 0) {
                    std::cerr << "Failed to receive image data" << std::endl;
                    break;
                }
                totalBytesReceived += bytesReceived;
            }

            if (totalBytesReceived == DATA_SIZE) {
                saveImage(imageData, DATA_SIZE);
            }
        }

    }
    close(serverSocket);
    return 0;
}