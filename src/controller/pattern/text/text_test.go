package text

import (
	"github.com/fogleman/gg"
	"github.com/thiswillgowell/light-controller/src/controller/pattern/text/face"
	"github.com/thiswillgowell/light-controller/src/display"
	"github.com/thiswillgowell/light-controller/src/piportal"
	"golang.org/x/image/colornames"
	"testing"
	"time"
)

func TestText(t *testing.T) {

	dc := gg.NewContextForImage(piportal.Fireplace.Image())
	dc.SetFontFace(face.LedSimpleSt(20))
	dc.SetColor(colornames.White)

	dc.DrawStringAnchored("Hello world", float64(dc.Image().Bounds().Max.X)/2, 40, 0.5, 0.5)
	dc.Stroke()

	display.DrawAndUpdate(display.NewRotation(piportal.Fireplace, display.MirrorAcrossY), dc.Image())
	<-time.After(time.Second)
}
