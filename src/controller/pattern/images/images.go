package images

import (
	"bytes"
	_ "embed"
	"image"
	"image/png"
)

//go:embed files/SuperMarioWorld3-3.png
var marioWorld []byte

func MarioWorld() image.Image {
	img, err := png.Decode(bytes.NewBuffer(marioWorld))
	if err != nil {
		panic(err)
	}
	return img
}
