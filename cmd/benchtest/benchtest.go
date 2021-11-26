package main

import (
	"context"
	"time"

	"github.com/thiswillgowell/light-controller/src/piportal/tx"
)

func main() {
	c, err := tx.NewUDPClient(context.Background(), "192.168.1.177:2390")

	if err != nil {
		panic(err)
	}
	buff := make([]byte, 1024)
	for i := 0; i < 100; i++ {
		err := c.WriteFrame(buff)
		if err != nil {
			panic(err)
		}
		<-time.After(time.Millisecond * 200)
	}
}
