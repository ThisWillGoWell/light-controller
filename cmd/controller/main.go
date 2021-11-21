package main

import (
	"github.com/thiswillgowell/light-controller/src/controller/pattern"
	"github.com/thiswillgowell/light-controller/src/daisy/daisy"
	"github.com/thiswillgowell/light-controller/src/piportal"
)

func main() {
	daisyDevice, err := daisy.InitDaisy()
	if err != nil {
		panic(err)
	}

	matrixDisplay, err := piportal.NewMatrix("192.168.1.63:8081")
	if err != nil {
		panic(err)
	}

	pattern.CenterVUBar(daisyDevice, matrixDisplay)
}
