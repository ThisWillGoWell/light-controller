package face

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)
import _ "embed"

//go:embed "BellotaText-Regular.ttf"
var ledSimpleStFile []byte

//go:embed "MonospaceBold.ttf"
var monospaceBold []byte

//go:embed "Monospace.ttf"
var monospace []byte

func parse(ttf []byte, points float64) font.Face {
	f, err := truetype.Parse(ttf)
	if err != nil {
		panic(err)
	}
	face := truetype.NewFace(f, &truetype.Options{
		Size: points,
		// Hinting: font.HintingFull,
	})
	return face
}

func LedSimpleSt(points float64) font.Face {
	return parse(ledSimpleStFile, points)
}

func MonospaceBold(points float64) font.Face {
	return parse(monospaceBold, points)
}

func Monospace(points float64) font.Face {
	return parse(monospace, points)
}
