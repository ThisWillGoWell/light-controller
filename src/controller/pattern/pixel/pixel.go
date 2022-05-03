package pixel

import (
	"github.com/fogleman/gg"
)

func IsEmpty(graphic *gg.Context, x, y int) bool {
	r, g, b, _ := graphic.Image().At(x, y).RGBA()
	return r == 0 && g == 0 && b == 0
}

func Copy(src, dest *gg.Context, x, y int) {
	r, g, b, a := src.Image().At(x, y).RGBA()
	dest.SetRGBA255(int(r), int(g), int(b), int(a))
	dest.SetPixel(x, y)
}
