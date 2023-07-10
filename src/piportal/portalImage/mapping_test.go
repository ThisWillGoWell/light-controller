package portalImage

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/image/colornames"
	"testing"
)

func TestThreeByTwoMapping(t *testing.T) {
	img := twoByThreeVerticalMappedPanel()
	mappedImg := img.MappedImage()

	img.Set(0, 0, colornames.Red)
	assert.Equal(t, colornames.Red, img.At(0, 0))
	assert.Equal(t, colornames.Red, mappedImg.At(0, 0))

	img.Set(0, 128, colornames.Purple)
	assert.Equal(t, colornames.Purple, mappedImg.At(64, 0))

	img.Set(63, 191, colornames.Orangered)
	assert.Equal(t, colornames.Orangered, mappedImg.At(64, 64))
}
