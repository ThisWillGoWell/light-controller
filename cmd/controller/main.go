package main

import (
	"github.com/thiswillgowell/light-controller/src/controller/pattern"
	"github.com/thiswillgowell/light-controller/src/daisy/daisy"
	"github.com/thiswillgowell/light-controller/src/display"
	"github.com/thiswillgowell/light-controller/src/piportal"
)

func main() {
	daisyDevice, err := daisy.InitDaisy()
	if err != nil {
		panic(err)
	}

	matrixRightDisplay, err := piportal.NewMatrix("192.168.1.53:8080")
	if err != nil {
		panic(err)
	}

	matrixLeftDisplay, err := piportal.NewMatrix("192.168.1.63:8080")
	if err != nil {
		panic(err)
	}

	combindedDispaly := display.NewMultiDisplay(display.ArrangementHorizontal,
		display.NewRotation(matrixLeftDisplay, display.OneEighty),
		matrixRightDisplay,
	)

	pattern.CenterHollowVUBarDouble(daisyDevice, combindedDispaly, 1)
}
