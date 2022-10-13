package music

import (
	"github.com/fogleman/gg"
	"github.com/thiswillgowell/light-controller/src/audio/specturm"
	"github.com/thiswillgowell/light-controller/src/display"
	"image/color"
)

type Params struct {
	BarWidth    int
	Display     display.Display
	InputMapper []int
	ColorGen    func(int, int) color.Color
	Channel     specturm.Channel
}

type CenterHollowVuMeter struct {
	dc     gg.Context
	params Params
}

func NewCenterHollowVuMeter(d display.Display, params Params) {
	targetImage := d.Image()
	numBars := targetImage.Bounds().Size().X / params.BarWidth
}

func (c *CenterHollowVuMeter) ProcessSpectrum(spectrum specturm.FrequencySpectrum) error {
	specturm.BinHeight()
	c.dc.SetColor(color.Transparent)
	c.dc.Clear()

	c.params.Display.Update()

}
