package text

import (
	"github.com/fogleman/gg"
	"github.com/thiswillgowell/light-controller/color"
	"github.com/thiswillgowell/light-controller/src/controller/pattern/text/face"
	"github.com/thiswillgowell/light-controller/src/display"
	"github.com/thiswillgowell/light-controller/src/piportal"
	"golang.org/x/image/colornames"
	"testing"
	"time"
)

func TestText(t *testing.T) {

	dc := gg.NewContextForImage(piportal.Fireplace.Image())
	largeFontFace := face.MonospaceBold(20)
	smallFontFace := face.Monospace(15)
	middle := gg.Point{
		X: float64(dc.Image().Bounds().Max.X) / 2.0,
		Y: float64(dc.Image().Bounds().Max.Y) / 2.0,
	}
	tick := time.NewTicker(time.Second / 100)
	for {
		<-tick.C
		dc.SetColor(color.Off)
		dc.Clear()
		dc.SetColor(colornames.Blueviolet)
		dc.SetFontFace(largeFontFace)
		currentTime := time.Now().Format("3:04:05.000")
		currentTime = currentTime[0 : len(currentTime)-2]
		_, firstHeight := dc.MeasureString(currentTime)
		dc.DrawStringAnchored(currentTime, middle.X, 20, 0.5, 0.5)

		dc.SetFontFace(smallFontFace)
		currentDate := time.Now().Format("Monday Jan 02")
		dc.DrawStringAnchored(currentDate, middle.X, 20+firstHeight+3, 0.5, 0.5)
		//dc.DrawStringAnchored("Hello world", float64(dc.Image().Bounds().Max.X)/2, 40, 0.5, 0.5)
		dc.Stroke()
		display.DrawAndUpdate(piportal.Fireplace, dc.Image())

	}

}
