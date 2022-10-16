package music

import (
	"github.com/fogleman/gg"
	"github.com/thiswillgowell/light-controller/src/audio/specturm"
	"github.com/thiswillgowell/light-controller/src/display"
	"image/color"
)

type Params struct {
	BarWidth int
	Display  display.Display
	//InputMapper []int
	// function that accepts int bar number, and bar height
	// then returns the []color.Color with length
	BarColors func(int, int) []color.Color
	Channel   specturm.Channel
}

type VuMeter struct {
	dc     *gg.Context
	params Params

	displayHeight int
	displayWidth  int

	numBars int
}

func NewVuMeter(params Params) *VuMeter {
	return &VuMeter{
		dc:            gg.NewContextForImage(params.Display.Image()),
		params:        params,
		displayHeight: params.Display.Image().Bounds().Size().Y,
		displayWidth:  params.Display.Image().Bounds().Size().X,
		numBars:       params.Display.Image().Bounds().Size().X / params.BarWidth,
	}
}

func (meter *VuMeter) ProcessSpectrum(s specturm.FrequencySpectrum) {
	barValues := s.Bin(meter.params.Channel, specturm.BinInput{
		MaxOutputValue: meter.displayHeight,
		NumOutputs:     meter.displayWidth,
		Interpolate:    true,
	}, true)

	meter.dc.SetColor(color.Transparent)
	meter.dc.Clear()
	for col, barHeight := range barValues {
		colors := meter.params.BarColors(col, barHeight)
		for i, c := range colors {
			meter.dc.SetColor(c)
			for j := 0; j < meter.params.BarWidth; j++ {
				meter.dc.SetPixel(col*meter.params.BarWidth+j, i)
			}
		}
	}
	display.DrawAndUpdate(meter.params.Display, meter.dc.Image())
}
