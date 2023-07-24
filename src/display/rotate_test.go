package display

import (
	"github.com/fogleman/gg"
	"github.com/stretchr/testify/assert"
	"golang.org/x/image/colornames"
	"testing"
)

func TestMirrorAcrossX(t *testing.T) {
	maxX := float64(32)
	maxY := float64(64)
	d := NewTestRGBA(32, 64, "mirror_x")
	mirrored := NewRotation(d, MirrorAcrossY)

	dc := gg.NewContextForImage(d.Image())
	dc.Clear()

	// fill with a line from 0,0 -> maxX,maxY
	dc.SetColor(colornames.Black)
	dc.DrawLine(0, 0, maxX, maxY)
	dc.Stroke()

	dc.SetColor(colornames.Purple)
	dc.DrawLine(0, 0, maxX-1, 0)
	dc.Stroke()

	//dc.SetColor(colornames.Red)
	//dc.DrawLine(maxX-1, 0, maxX-1, maxY-1)
	//dc.Stroke()
	//
	//dc.SetColor(colornames.Blue)
	//dc.DrawLine(maxX-1, maxY-1, 0, maxY-1)
	//dc.Stroke()

	dc.SetColor(colornames.Gold)
	dc.DrawLine(0, maxY-1, 0, 0)
	dc.Stroke()

	DrawAndUpdate(mirrored, dc.Image())

	// now make sure the images are the same but mirrored
	for x := 0; x < int(maxX); x++ {
		for y := 0; y < int(maxY); y++ {
			assert.Equal(t, d.Image().At(x, y), mirrored.Image().At(int(maxX)-x-1, y))
		}
	}

}
