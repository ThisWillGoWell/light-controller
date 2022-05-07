package face

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)
import _ "embed"

//go:embed "BellotaText-Regular.ttf"
var ledSimpleStFile []byte

func LedSimpleSt(points float64) font.Face {
	f, err := truetype.Parse(ledSimpleStFile)
	if err != nil {
		panic(err)
	}
	face := truetype.NewFace(f, &truetype.Options{
		Size: points,
		// Hinting: font.HintingFull,
	})
	return face
}
