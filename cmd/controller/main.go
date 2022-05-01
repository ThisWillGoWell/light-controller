package main

import (
	"github.com/thiswillgowell/light-controller/src/controller/pattern/music"
	"github.com/thiswillgowell/light-controller/src/daisy/daisy"
	"github.com/thiswillgowell/light-controller/src/piportal"
)

func main() {
	daisyDevice, err := daisy.InitDaisy()
	if err != nil {
		panic(err)
	}
	p1, err := piportal.NewMatrix("192.168.1.53:8080", piportal.Right)
	if err != nil {
		panic(err)
	}
	//p2, err := piportal.NewMatrix("192.168.1.83:8080", piportal.Left)
	//if err != nil {
	//	panic(err)
	//}

	//p := display.NewRotation(display.NewMultiDisplay(display.ArrangementVertical, display.NewRotation(p2, display.MirrorAcrossY), p1), display.MirrorAcrossY)

	music.CenterHollowVUBarDouble(daisyDevice, p1, 1)
}
