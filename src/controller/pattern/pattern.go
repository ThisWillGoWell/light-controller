package pattern

import (
	"github.com/thiswillgowell/light-controller/color"
	"github.com/thiswillgowell/light-controller/src/controller"
	"github.com/thiswillgowell/light-controller/src/controller/pattern/music"
	"github.com/thiswillgowell/light-controller/src/daisy/daisy"
)

type Pattern interface {
	GetNextValue() [][]color.Color
}

type MusicPattern struct {
}

func CenterVUBar(daisyDevice *daisy.Daisy, display controller.Display) {

	lut := music.BuildIndexLUT(music.MappingInput{
		InputSize:         daisy.NumFrequencies,
		OutputSize:        display.Width(),
		InputPercentages:  []float64{1},
		OutputPercentages: []float64{1},
	})

	for channel := range daisyDevice.FFTChannel {
		controller.ForEach(display, controller.DarkenDisplay(20))
		leftBars := music.BinHeight(music.BinInput{
			Input:          channel[0],
			MaxInputValue:  100.0,
			MaxOutputValue: display.Height() / 2,
			BinningLut:     lut,
			NumOutputs:     display.Width(),
		})
		// left is top half of the display
		startRow := display.Height() / 2
		for col, height := range leftBars {
			for i := 0; i < height; i++ {
				display.SetPixel(startRow+i, col, color.Darkred)
			}
		}

		rightBars := music.BinHeight(music.BinInput{
			Input:          channel[1],
			MaxInputValue:  100.0,
			MaxOutputValue: display.Height() / 2,
			BinningLut:     lut,
			NumOutputs:     display.Width(),
		})

		// right is top half of the display, but looking down
		startRow = display.Height()/2 - 1
		for col, height := range rightBars {
			for i := 0; i < height; i++ {
				display.SetPixel(startRow-i, col, color.Darkred)
			}
		}
		display.Send()
	}
}
