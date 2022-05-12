package nes

import (
	_ "embed"
	"github.com/fogleman/nes/nes"
	"github.com/thiswillgowell/light-controller/src/display"
	"golang.org/x/image/draw"
)

func drawFrame(d display.Display, c *nes.Console) {
	draw.BiLinear.Scale(d.Image(), d.Image().Bounds(), c.Buffer(), c.Buffer().Bounds(), draw.Over, nil)
	d.Update()
}
