import argparse
import time
import sys
import os
import socket
from rgbmatrix import RGBMatrix, RGBMatrixOptions

localIP = "0.0.0.0"
localPort = 8081
bufferSize = 3 * 64 * 64


# sys.path.append(os.path.abspath(os.path.dirname(__file__) + '/..'))


class PixelPortal(object):
    def __init__(self, *args, **kwargs):
        print("starting server")
        self.UDPServerSocket = socket.socket(family=socket.AF_INET, type=socket.SOCK_DGRAM)
        self.UDPServerSocket.bind((localIP, localPort))
        options = RGBMatrixOptions()
        options.rows = 32
        options.cols = 64
        options.chain_length = 2
        options.parallel = 1
        options.row_address_type = 0
        options.multiplexing = 0
        options.pwm_bits = 11
        options.brightness = 100
        options.pwm_lsb_nanoseconds = 130
        options.led_rgb_sequence = "RGB"
        options.pixel_mapper_config = "U-mapper"
        options.show_refresh_rate = 0
        options.gpio_slowdown = 4
        self.matrix = RGBMatrix(options=options)
        return

    def run(self):
        canvas = self.matrix.CreateFrameCanvas()
        while True:
            bytesAddressPair = self.UDPServerSocket.recvfrom(bufferSize)
            message = bytesAddressPair[0]
            r, g, b = 0, 0, 0
            row = message[0]
            for i in range(0, len(message)):
                if i % 3 == 0:
                    r = int(message[i])
                elif i % 3 == 1:
                    g = int(message[i])
                elif i % 3 == 2:
                    b = int(message[i])
                    canvas.SetPixel(i // 3, row, r, g, b)

            if row % (self.matrix.height - 1) == 0 and row != 0:
                canvas = self.matrix.SwapOnVSync(canvas)

    def start(self):
        try:
            # Start loop
            print("Press CTRL-C to stop sample")
            self.run()
        except KeyboardInterrupt:
            print("Exiting\n")
            sys.exit(0)

        return True


if __name__ == "__main__":
    PixelPortal().start()
