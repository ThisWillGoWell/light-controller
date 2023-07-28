package graphics

import (
	"image"
	"image/color"
)

type TextWriter func(message string, c color.Color, position image.Point)
