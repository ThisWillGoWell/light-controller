package display

import (
	"github.com/fogleman/gg"
	"go.uber.org/zap"
	"golang.org/x/image/draw"
	"image"
	"strings"
)

type Display interface {
	Image() draw.Image
	Update()
}

func DrawAndUpdate(dest Display, src image.Image) {
	draw.Draw(dest.Image(), dest.Image().Bounds(), src, image.Point{}, draw.Src)
	dest.Update()
}

func NewTestRGBA(x, y int, filename string) Display {
	return &rgba{
		img:      image.NewRGBA(image.Rect(0, 0, x, y)),
		fileName: strings.TrimRight(filename, ".png") + ".png",
	}
}

func NewRGBA(x, y int) Display {
	return &rgba{
		img: image.NewRGBA(image.Rect(0, 0, x, y)),
	}
}

type rgba struct {
	img      *image.RGBA
	fileName string
}

func (R *rgba) Image() draw.Image {
	return R.img
}

func (R *rgba) Update() {
	if R.fileName != "" {
		if err := gg.NewContextForImage(R.img).SavePNG(R.fileName); err != nil {
			zap.L().Error("failed to update", zap.Error(err))
		}
	}
}
