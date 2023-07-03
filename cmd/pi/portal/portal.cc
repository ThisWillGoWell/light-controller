// -*- mode: c++; c-basic-offset: 2; indent-tabs-mode: nil; -*-
// Small example how to use the library.
// For more examples, look at demo-main.cc
//
// This code is public domain
// (but note, that the led-matrix library this depends on is GPL v2)

#include "led-matrix.h"
//#include "graphics.h"
// led matrix
#include <unistd.h>
#include <math.h>
#include <stdio.h>
#include <signal.h>

// tcp server

#include <stdlib.h>
#include <unistd.h>
#include <errno.h>
#include <string.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <netdb.h>
#include <arpa/inet.h>
#include <sys/wait.h>


using rgb_matrix::RGBMatrix;
using rgb_matrix::Canvas;
using rgb_matrix::FrameCanvas;

#define PORT "8080"
#define MAXDATASIZE 193 // max number of bytes we can get at once

void sigchld_handler(int s)
{
    // waitpid() might overwrite errno, so we save and restore it:
    int saved_errno = errno;

    while(waitpid(-1, NULL, WNOHANG) > 0);

    errno = saved_errno;
}

// get sockaddr, IPv4 or IPv6:
void *get_in_addr(struct sockaddr *sa)
{
    if (sa->sa_family == AF_INET) {
        return &(((struct sockaddr_in*)sa)->sin_addr);
    }

    return &(((struct sockaddr_in6*)sa)->sin6_addr);
}

int row;
// process
bool processPacket( int byteCount, char buf[], FrameCanvas *offscreen){
  row = int(buf[0]);
  printf("row: %d\n", row);
   for(int led=0; led <(byteCount-1)/3; led++) {
    printf("row: %d, col: %d, r: %d, g: %d, b:%d\n", row, led, int(buf[1 + led*3]), int(buf[2 + led*3]), int(buf[3 + led*3]));
    //offscreen->SetPixel(led, row, int(buf[1 + led*3]), int(buf[2 + led*3]), int(buf[3 + led*3]));
    offscreen->SetPixel(led, row, 127, 20, 24);

   }
   return row == 63;
 }

int main(int argc, char *argv[]) {

  RGBMatrix::Options defaults;
  rgb_matrix::RuntimeOptions runtime_opt;

  runtime_opt.gpio_slowdown = 5;
  defaults.rows = 32;
  defaults.cols = 128;
  defaults.chain_length = 2;
  defaults.parallel = 3;
  defaults.show_refresh_rate = false;
  RGBMatrix *matrix = RGBMatrix::CreateFromOptions(defaults, runtime_opt);
  FrameCanvas *offscreen = matrix->CreateFrameCanvas();

  int sockfd, new_fd;  // listen on sock_fd, new connection on new_fd
  struct addrinfo hints, *servinfo, *p;
  struct sockaddr_storage their_addr; // connector's address information
  socklen_t sin_size;
  struct sigaction sa;
  int yes=1;
  char s[INET6_ADDRSTRLEN];
  int rv;

  memset(&hints, 0, sizeof hints);
  hints.ai_family = AF_UNSPEC;
  hints.ai_socktype = SOCK_STREAM;
  hints.ai_flags = AI_PASSIVE; // use my IP


  char buf[MAXDATASIZE];
  int bytecount;
  if ((rv = getaddrinfo(NULL, PORT, &hints, &servinfo)) != 0) {
     fprintf(stderr, "getaddrinfo: %s\n", gai_strerror(rv));
     return 1;
  }

  // loop through all the results and bind to the first we can
  for(p = servinfo; p != NULL; p = p->ai_next) {
     if ((sockfd = socket(p->ai_family, p->ai_socktype,
             p->ai_protocol)) == -1) {
         perror("server: socket");
         continue;
     }

     if (setsockopt(sockfd, SOL_SOCKET, SO_REUSEADDR, &yes,
             sizeof(int)) == -1) {
         perror("setsockopt");
         exit(1);
     }

     if (bind(sockfd, p->ai_addr, p->ai_addrlen) == -1) {
         close(sockfd);
         perror("server: bind");
         continue;
     }

     break;
  }

  freeaddrinfo(servinfo); // all done with this structure

  if (p == NULL)  {
     fprintf(stderr, "server: failed to bind\n");
     exit(1);
  }

  if (listen(sockfd, 5) == -1) {
     perror("listen");
     exit(1);
  }

  sa.sa_handler = sigchld_handler; // reap all dead processes
  sigemptyset(&sa.sa_mask);
  sa.sa_flags = SA_RESTART;
  if (sigaction(SIGCHLD, &sa, NULL) == -1) {
     perror("sigaction");
     exit(1);
  }
  printf("server: waiting for connections...\n");

  while(1) {  // main accept() loop
     sin_size = sizeof their_addr;
     new_fd = accept(sockfd, (struct sockaddr *)&their_addr, &sin_size);
     if (new_fd == -1) {
         perror("accept");
         continue;
     }

   inet_ntop(their_addr.ss_family,
       get_in_addr((struct sockaddr *)&their_addr),
       s, sizeof s);
   printf("server: got connection from %s\n", s);


    bool running = true;
    while(running){
       bytecount = recv(new_fd, buf, MAXDATASIZE, 0);
       if(bytecount < 1){
          printf("done\n");
          running = false;
       } else {
          row = int(buf[0]);
          for(int led=0; led <(bytecount-1)/3; led++) {
            int newLed = led;
            int newRow = row;
            if(row < 32) {
              newRow = 31 - row;
              newLed = 63 - led;
            } else if(row >= 64) {
               newRow = 64 + (95 - row);
               newLed = 63 - led;
            }

            printf("len: %d, row: %d, col: %d, r:%d, g:%d, b:%d\n", bytecount, newRow, newLed,  int(buf[1 + led*3]), int(buf[2 + led*3]), int(buf[3 + led*3]));
            offscreen->SetPixel(newLed, newRow, int(buf[1 + led*3]), int(buf[2 + led*3]), int(buf[3 + led*3]));
          }
          if(int(buf[0]) == 95){
            offscreen = matrix->SwapOnVSync(offscreen);
          }
      }
    }

   printf("%s\n", "closed!");
   close(new_fd);  // parent doesn't need this
  }

   return 0;
}
