package main

import (
	"github.com/thiswillgowell/light-controller/src/controller/pattern/music"
	"github.com/thiswillgowell/light-controller/src/daisy/daisy"
	"github.com/thiswillgowell/light-controller/src/display"
	"github.com/thiswillgowell/light-controller/src/piportal"
)

func main() {
	daisyDevice, err := daisy.Init()
	if err != nil {
		panic(err)
	}
	p1, err := piportal.NewMatrix("192.168.1.53:8080", piportal.Right)
	if err != nil {
		panic(err)
	}
	p2, err := piportal.NewMatrix("192.168.1.83:8080", piportal.Left)
	if err != nil {
		panic(err)
	}

	p := display.NewRotation(display.NewMultiDisplay(display.ArrangementVertical, display.NewRotation(p2, display.MirrorAcrossY), p1), display.CounterClockwise)
	//p := display.NewMirrorDisplay(p1, p2)
	music.CenterHollowVUBarDouble(daisyDevice, p, 2)
}
