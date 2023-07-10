package portalImage

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/image/colornames"
	"image/png"
	"os"
	"testing"
)

func TestThreeByTwoMapping(t *testing.T) {
	img := twoByThreeVerticalMappedPanel()
	mappedImg := img.MappedImage()

	img.Set(0, 0, colornames.Blue)
	img.Set(0, 127, colornames.Red)
	img.Set(31, 0, colornames.Red)
	img.Set(31, 127, colornames.Green)
	assert.Equal(t, colornames.Blue, img.At(0, 0))
	assert.Equal(t, colornames.Blue, mappedImg.At(127, 0))

	f, err := os.Create("test.png")
	if err != nil {
		panic(err)
	}
	if err := png.Encode(f, mappedImg); err != nil {
		panic(err)
	}

	img.Set(0, 128, colornames.Purple)
	assert.Equal(t, colornames.Purple, mappedImg.At(64, 0))

	img.Set(63, 191, colornames.Orangered)
	assert.Equal(t, colornames.Orangered, mappedImg.At(64, 64))
}
