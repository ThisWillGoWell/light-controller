package main

import (
	"context"
	"github.com/thiswillgowell/light-controller/src/piportal/tx"
	"time"
)

func main(){
	c, err := tx.NewClient(context.Background(), "192.168.1.177:2390")

	if err != nil {
		panic(err)
	}
	buff := make([]byte, 1024)
	for i := 0;i<100;i++ {
		err := c.WriteFrame(buff)
		if err != nil{
			panic(err)
		}
		<- time.After(time.Millisecond * 200)
	}
}