package display

import (
	"golang.org/x/image/draw"
	"image"
)

type Display interface {
	Image() draw.Image
	Update()
}

func DrawAndUpdate(dest Display, src image.Image) {
	draw.Draw(dest.Image(), dest.Image().Bounds(), src, image.Point{}, draw.Src)
	dest.Update()
}

func NewRGBA(x, y int) Display {
	return &rgba{
		img: image.NewRGBA(image.Rect(0, 0, x, y)),
	}
}

type rgba struct {
	img *image.RGBA
}

func (R *rgba) Image() draw.Image {
	return R.img
}

func (R *rgba) Update() {}
