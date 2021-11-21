package controller

import (
	"github.com/thiswillgowell/light-controller/color"
	"github.com/thiswillgowell/light-controller/src/daisy/daisy"
)

type Controller struct {
	Daisy    daisy.Daisy
	Displays map[string]Display
}

type Display interface {
	Height() int
	Width() int
	Clear()
	GetPixel(row, col int) color.Color
	SetPixel(row, col int, c color.Color)
	Send() error
}

func ForEach(d Display, each func(row, col int, c color.Color) color.Color) {
	for r := 0; r < d.Height(); r++ {
		for c := 0; c < d.Width(); c++ {
			d.SetPixel(r, c, each(r, c, d.GetPixel(r, c)))
		}
	}
}

func DarkenDisplay(amount float64) func(int, int, color.Color) color.Color {
	return func(_, _ int, c color.Color) color.Color {
		var newV uint8
		if c.V <= 10 {
			newV = 0
		} else {
			newV = uint8(float64(newV) * amount)
		}
		return color.FromHsv(c.H, c.S, newV)
	}
}
