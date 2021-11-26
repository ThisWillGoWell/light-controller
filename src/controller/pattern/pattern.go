package pattern

import (
	"github.com/thiswillgowell/light-controller/color"
	"github.com/thiswillgowell/light-controller/src/controller"
	"github.com/thiswillgowell/light-controller/src/controller/pattern/music"
	"github.com/thiswillgowell/light-controller/src/daisy/daisy"
	"github.com/thiswillgowell/light-controller/src/display"
)

type Pattern interface {
	GetNextValue() [][]color.Color
}

type MusicPattern struct {
}

func CenterVUBar(daisyDevice *daisy.Daisy, display display.Display) {
	colors := color.LinerGradient(color.Darkred, color.Purple, display.Cols(), false)
	lut := music.BuildIndexLUT(music.MappingInput{
		InputSize:         daisy.NumFrequencies,
		OutputSize:        display.Cols(),
		InputPercentages:  []float64{0.4, 0.3, 0.3},
		OutputPercentages: []float64{0.5, 0.4, 0.1},
		Reversed:          true,
	})

	for channel := range daisyDevice.FFTChannel {
		controller.ForEach(display, controller.DarkenDisplay(.3))
		leftBars := music.BinHeight(music.BinInput{
			Input:          channel[0],
			MaxInputValue:  200.0,
			MaxOutputValue: display.Rows() / 2,
			BinningLut:     lut,
			NumOutputs:     display.Cols(),
		})
		// left is top half of the display
		startRow := display.Rows() / 2
		for col, height := range leftBars {
			for i := 0; i < height; i++ {
				display.SetPixel(startRow+i, col, colors[col])
			}
		}

		rightBars := music.BinHeight(music.BinInput{
			Input:          channel[1],
			MaxInputValue:  200.0,
			MaxOutputValue: display.Rows() / 2,
			BinningLut:     lut,
			NumOutputs:     display.Cols(),
		})

		// right is top half of the display, but looking down
		startRow = display.Rows()/2 - 1
		for col, height := range rightBars {
			for i := 0; i < height; i++ {
				display.SetPixel(startRow-i, col, colors[col])
			}
		}
		display.Send()
	}
}

func CenterHollowVUBar(daisyDevice *daisy.Daisy, display display.Display) {
	colors := color.LinerGradient(color.Darkred, color.Purple, display.Cols(), false)
	lut := music.BuildIndexLUT(music.MappingInput{
		InputSize:         daisy.NumFrequencies,
		OutputSize:        display.Cols(),
		InputPercentages:  []float64{0.4, 0.3, 0.3},
		OutputPercentages: []float64{0.5, 0.4, 0.1},
		Reversed:          true,
	})

	for channel := range daisyDevice.FFTChannel {
		controller.ForEach(display, controller.DarkenDisplay(.2))
		//display.Clear()
		leftBars := music.BinHeight(music.BinInput{
			Input:          channel[0],
			MaxInputValue:  200.0,
			MaxOutputValue: display.Rows() / 2,
			BinningLut:     lut,
			NumOutputs:     display.Cols(),
		})
		// left is top half of the display
		startRow := display.Rows() / 2
		for col, height := range leftBars {
			for i := 0; i < height; i++ {
				c := colors[col].DarkenPercentage(float64(i+2) / float64(height))
				display.SetPixel(startRow+i, col, c)
			}
		}

		rightBars := music.BinHeight(music.BinInput{
			Input:          channel[1],
			MaxInputValue:  200.0,
			MaxOutputValue: display.Rows() / 2,
			BinningLut:     lut,
			NumOutputs:     display.Cols(),
		})

		// right is top half of the display, but looking down
		startRow = display.Rows()/2 - 1
		for col, height := range rightBars {
			for i := 0; i < height; i++ {
				c := colors[col].DarkenPercentage(float64(i+2) / float64(height))
				display.SetPixel(startRow-i, col, c)
			}
		}
		display.Send()
	}
}

func CenterHollowVUBarDouble(daisyDevice *daisy.Daisy, display display.Display, barWidth int) {

	numBars := display.Cols() / barWidth
	colors := color.LinerGradient(color.Darkred, color.Purple, numBars, true)
	lut := music.BuildIndexLUT(music.MappingInput{
		InputSize:         daisy.NumFrequencies,
		OutputSize:        numBars,
		InputPercentages:  []float64{0.4, 0.3, 0.3},
		OutputPercentages: []float64{0.5, 0.4, 0.1},
		Reversed:          false,
	})

	for channel := range daisyDevice.FFTChannel {
		controller.ForEach(display, controller.DarkenDisplay(.3))
		leftBars := removeDeadZones(music.BinHeight(music.BinInput{
			Input:          channel[0],
			MaxInputValue:  200.0,
			MaxOutputValue: display.Rows() / 2,
			BinningLut:     lut,
			NumOutputs:     numBars,
			Interpolate:    true,
		}))
		// left is top half of the display
		startRow := display.Rows() / 2
		for col, height := range leftBars {
			c := colors[col]
			if height < 3 {
				c = c.DarkenPercentage(float64(4-height) * .1)
			}

			for i := 0; i < height; i++ {
				c := c.DarkenPercentage(float64(i+2) / float64(height))
				for j := 0; j < barWidth; j++ {
					display.SetPixel(startRow+i, col*barWidth+j, c)
				}
			}
		}

		rightBars := removeDeadZones(music.BinHeight(music.BinInput{
			Input:          channel[1],
			MaxInputValue:  200.0,
			MaxOutputValue: display.Rows() / 2,
			BinningLut:     lut,
			NumOutputs:     numBars,
			Interpolate:    true,
		}))

		// right is top half of the display, but looking down
		startRow = display.Rows()/2 - 1
		for col, height := range rightBars {
			c := colors[col]
			if height < 3 {
				c = c.DarkenPercentage(float64(4-height) * .1)
			}

			for i := 0; i < height; i++ {
				c := c.DarkenPercentage(float64(i+2) / float64(height))
				for j := 0; j < barWidth; j++ {
					display.SetPixel(startRow-i, col*barWidth+j, c)
				}
			}
		}
		display.Send()
	}
}

func FallingVuMeter(daisyDevice *daisy.Daisy, display display.Display) {
	colors := color.LinerGradient(color.Darkred, color.ForestGreen, display.Rows(), true)
	lut := music.BuildIndexLUT(music.MappingInput{
		InputSize:         daisy.NumFrequencies,
		OutputSize:        display.Cols() / 2,
		InputPercentages:  []float64{0.4, 0.3, 0.3},
		OutputPercentages: []float64{0.5, 0.4, 0.1},
		Reversed:          false,
	})

	for channel := range daisyDevice.FFTChannel {
		controller.ForEach(display, controller.DarkenDisplay(.3))
		leftBars := music.BinHeight(music.BinInput{
			Input:          channel[0],
			MaxInputValue:  200.0,
			MaxOutputValue: display.Rows(),
			BinningLut:     lut,
			NumOutputs:     display.Cols() / 2,
		})
		// left is top half of the display
		startRow := display.Rows() / 2
		for col, height := range leftBars {
			c := colors[col]
			if height < 4 {
				c = c.DarkenPercentage(float64(5-height) * .1)
			}

			for i := 0; i < height; i++ {
				c := c.DarkenPercentage(float64(i) / float64(height))
				display.SetPixel(startRow+i, col*2, c)
				display.SetPixel(startRow+i, col*2+1, c)
			}
		}

		rightBars := music.BinHeight(music.BinInput{
			Input:          channel[1],
			MaxInputValue:  200.0,
			MaxOutputValue: display.Rows() / 2,
			BinningLut:     lut,
			NumOutputs:     display.Cols() / 2,
		})

		// right is top half of the display, but looking down
		startRow = display.Rows()/2 - 1
		for col, height := range rightBars {
			c := colors[col]
			if height < 4 {
				c = c.DarkenPercentage(float64(5-height) * .1)
			}
			for i := 0; i < height; i++ {
				c := c.DarkenPercentage(float64(i) / float64(height))
				display.SetPixel(startRow-i, col*2, c)
				display.SetPixel(startRow-i, col*2+1, c)
			}
		}
		display.Send()
	}
}

func removeDeadZones(input []int) []int {
	for i, ele := range input {
		if i == 0 || i == len(input)-1 {
			continue
		}
		// leave alone if bigger than one
		if ele > 1 {
			continue
		}
		// set to zero if previous value < 1 and next value is < 1

		if input[i-1] <= 1 && input[i+1] <= 1 {
			input[i] = 0
		}
	}
	return input
}
