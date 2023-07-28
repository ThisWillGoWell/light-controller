package text

import (
	"fmt"
	"github.com/thiswillgowell/light-controller/src/display"
	"golang.org/x/image/colornames"
	"image"
	"testing"
)

func TestFonts(t *testing.T) {
	for fontType, f := range fonts {
		testImage := display.NewTestRGBA(64, 96, fmt.Sprintf("%d-%s", fontType, f.Name))
		WriteOnImage("abc 123", fontType, colornames.Black, image.Point{}, testImage.Image())
		testImage.Update()
	}
}
