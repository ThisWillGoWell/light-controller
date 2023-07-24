package portals

import (
	"github.com/fogleman/gg"
	"github.com/thiswillgowell/light-controller/src/display"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"testing"
	"time"
)

func drawTestImage(p display.Display) {
	img := image.NewRGBA(p.Image().Bounds())
	dc := gg.NewContextForRGBA(img)
	maxX, maxY := float64(p.Image().Bounds().Max.X), float64(p.Image().Bounds().Max.Y)
	for x := 0; x < p.Image().Bounds().Max.X; x++ {
		dc.Clear()
		grad := gg.NewRadialGradient(maxX/2, maxY/2, 10, maxX/2, maxY/2, 80)
		grad.AddColorStop(0, color.RGBA{10, 50, 0, 255})
		grad.AddColorStop(1, color.RGBA{60, 0, 40, 255})

		dc.SetFillStyle(grad)
		dc.DrawRectangle(0, 0, maxX, maxY)
		dc.Fill()

		dc.SetColor(color.White)
		dc.DrawCircle(float64(x), float64(x)/float64(p.Image().Bounds().Max.X)*float64(p.Image().Bounds().Max.Y), 10)
		dc.Stroke()
		dc.DrawCircle(float64(x), 0, 80)
		dc.Stroke()
		dc.DrawLine(0, 0, float64(p.Image().Bounds().Max.X), float64(p.Image().Bounds().Max.Y))
		dc.Stroke()

		dc.SetColor(colornames.Purple)
		dc.DrawLine(0, 0, maxX-1, 0)
		dc.Stroke()
		dc.SetColor(colornames.Red)
		dc.DrawLine(maxX-1, 0, maxX-1, maxY-1)
		dc.Stroke()
		dc.SetColor(colornames.Blue)
		dc.DrawLine(maxX-1, maxY-1, 0, maxY-1)
		dc.Stroke()
		dc.SetColor(colornames.Gold)
		dc.DrawLine(0, maxY-1, 0, 0)
		dc.Stroke()
		display.DrawAndUpdate(p, dc.Image())
		<-time.After(time.Second)
	}
}

func TestBothVerticals(t *testing.T) {
	drawTestImage(BothVerticals)
}
