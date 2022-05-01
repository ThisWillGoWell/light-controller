package controller

import (
	"github.com/thiswillgowell/light-controller/color_2"
	"github.com/thiswillgowell/light-controller/src/daisy/daisy"
	"github.com/thiswillgowell/light-controller/src/display"
)

type Controller struct {
	Daisy    daisy.Daisy
	Displays map[string]display.Display
}

//func ForEach(d display.Display, each func(row, col int, c color_2.Color) color_2.Color) {
//	for r := 0; r < d.Image().Bounds()..Rows(); r++ {
//		for c := 0; c < d.Cols(); c++ {
//			d.SetPixel(r, c, each(r, c, d.GetPixel(r, c)))
//		}
//	}
//}

func DarkenDisplay(amount float64) func(int, int, color_2.Color) color_2.Color {
	return func(_, _ int, c color_2.Color) color_2.Color {
		var newV uint8
		if c.V <= 10 {
			newV = 0
		} else {
			newV = uint8(float64(newV) * amount)
		}
		return color_2.FromHsv(c.H, c.S, newV)
	}
}
