package color

import (
	"image/color"

	"golang.org/x/image/colornames"
)

var (
	Off           = fromGolangColor(color.RGBA{})
	Crimson       = fromGolangColor(colornames.Crimson)
	Darkmagenta   = fromGolangColor(colornames.Darkmagenta)
	Darkred       = fromGolangColor(colornames.Darkred)
	Deeppink      = fromGolangColor(colornames.Deeppink)
	Indigo        = fromGolangColor(colornames.Indigo)
	Mediumblue    = fromGolangColor(colornames.Mediumblue)
	Midnightblue  = fromGolangColor(colornames.Midnightblue)
	Firebrick     = fromGolangColor(colornames.Firebrick)
	Darkorange    = fromGolangColor(colornames.Darkorange)
	Orangered     = fromGolangColor(colornames.Orangered)
	Darkslateblue = fromGolangColor(colornames.Darkslateblue)
	Darkorchid    = fromGolangColor(colornames.Darkorchid)
	Maroon        = fromGolangColor(colornames.Maroon)
	Purple        = fromGolangColor(colornames.Purple)
	Yellow        = fromGolangColor(colornames.Yellow)
	ForestGreen   = fromGolangColor(colornames.Forestgreen)
)

func fromGolangColor(c color.Color) Color {
	r2, g2, b2, _ := c.RGBA()
	r, g, b := uint8(r2), uint8(g2), uint8(b2)
	h, s, v := RgbToHsv(r, g, b)
	return Color{
		H: h,
		S: s,
		V: v,
		R: r,
		G: g,
		B: b,
	}
}
