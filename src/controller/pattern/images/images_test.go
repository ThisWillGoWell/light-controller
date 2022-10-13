package images

import (
	"github.com/fogleman/gg"
	"github.com/thiswillgowell/light-controller/src/controller/pattern/text/face"
	"github.com/thiswillgowell/light-controller/src/display"
	"github.com/thiswillgowell/light-controller/src/live"
	"github.com/thiswillgowell/light-controller/src/piportal"
	"golang.org/x/image/colornames"
	"golang.org/x/image/draw"
	"image"
	"testing"
	"time"
)

func TestText(t *testing.T) {
	d := piportal.Fireplace
	subD := display.NewSubscription(d)
	go live.RunServer(subD)
	dc := gg.NewContextForImage(subD.Image())
	dc.SetFontFace(face.LedSimpleSt(20))
	dc.SetColor(colornames.White)
	marioWorld := MarioWorld()

	in := marioWorld.Bounds().Max
	out := d.Image().Bounds().Max

	scaleFactor := float64(in.Y) / float64(out.Y)
	xAmount := int(float64(out.X) * scaleFactor)
	ticker := time.NewTicker(time.Second / 75)
	for x := 0; x < in.X-xAmount; x++ {
		<-ticker.C
		draw.CatmullRom.Scale(subD.Image(), subD.Image().Bounds(), marioWorld, image.Rect(x, 0, x+xAmount, in.Y), draw.Over, nil)
		subD.Update()
	}
}
