package piportal

import (
	"github.com/thiswillgowell/light-controller/src/display"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"math/rand"
	"testing"
	"time"

	"github.com/fogleman/gg"
)

func random() float64 {
	return rand.Float64()*2 - 1
}

func point() (x, y float64) {
	return random(), random()
}

func drawCurve(dc *gg.Context) {
	dc.SetRGBA(0, 0, 0, 0.1)
	dc.FillPreserve()
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(12)
	dc.Stroke()
}

func drawPoints(dc *gg.Context) {
	dc.SetRGBA(1, 0, 0, 0.5)
	dc.SetLineWidth(2)
	dc.Stroke()
}

func randomQuadratic(dc *gg.Context) {
	x0, y0 := point()
	x1, y1 := point()
	x2, y2 := point()
	dc.MoveTo(x0, y0)
	dc.QuadraticTo(x1, y1, x2, y2)
	drawCurve(dc)
	dc.MoveTo(x0, y0)
	dc.LineTo(x1, y1)
	dc.LineTo(x2, y2)
	drawPoints(dc)
}

func randomCubic(dc *gg.Context) {
	x0, y0 := point()
	x1, y1 := point()
	x2, y2 := point()
	x3, y3 := point()
	dc.MoveTo(x0, y0)
	dc.CubicTo(x1, y1, x2, y2, x3, y3)
	drawCurve(dc)
	dc.MoveTo(x0, y0)
	dc.LineTo(x1, y1)
	dc.LineTo(x2, y2)
	dc.LineTo(x3, y3)
	drawPoints(dc)
}

func TestPortals(t *testing.T) {
	//	p1, err := NewMatrix("192.168.1.84:8080", TopLeft)

	//p := display.NewMirrorDisplay(Fireplace)
	p := display.NewMirrorDisplay(TopLeftDisplay, TopRightDisplay, BottomLeftDisplay, BottomRightDisplay)
	//p1, err := NewMatrix("192.168.1.53:8080", TopLeft)
	//if err != nil {
	//	panic(err)
	//}
	////	p2, err := NewMatrix("192.168.1.106:8080", TopRight)
	//p2, err := NewMatrix("192.168.1.83:8080", TopRight)
	//if err != nil {
	//	panic(err)
	//}
	//
	//p := display.NewRotation(display.NewMultiDisplay(display.ArrangementVertical, display.NewRotation(p2, display.MirrorAcrossY), p1), display.MirrorAcrossY)

	//p := display.NewMirrorDisplay(p1, display.NewRotation(p2, display.MirrorAcrossY))

	img := image.NewRGBA(p.Image().Bounds())
	dc := gg.NewContextForRGBA(img)
	maxX, maxY := float64(p.Image().Bounds().Max.X), float64(p.Image().Bounds().Max.Y)
	for x := 0; x < p.Image().Bounds().Max.X; x++ {
		dc.Clear()
		grad := gg.NewRadialGradient(maxX/2, maxY/2, 10, maxX/2, maxY/2, 80)
		grad.AddColorStop(0, color.RGBA{10, 50, 0, 255})
		grad.AddColorStop(1, color.RGBA{60, 0, 40, 255})

		dc.SetFillStyle(grad)
		dc.DrawRectangle(1, 1, maxX-2, maxY-2)
		dc.Fill()

		dc.SetColor(color.White)
		dc.DrawCircle(float64(x), float64(x)/float64(p.Image().Bounds().Max.X)*float64(p.Image().Bounds().Max.Y), 10)
		dc.Stroke()
		dc.DrawCircle(float64(x), 0, 80)
		dc.Stroke()
		dc.DrawLine(0, 0, float64(p.Image().Bounds().Max.X), float64(p.Image().Bounds().Max.Y))
		dc.Stroke()

		dc.SetColor(colornames.Hotpink)
		dc.DrawLine(0, 0, maxX, 1)
		dc.Stroke()
		display.DrawAndUpdate(p, dc.Image())
		<-time.After(time.Second)
	}

}

func TestPortal(t *testing.T) {

	p, err := NewMatrix("192.168.1.53:8080", TopLeft)
	if err != nil {
		panic(err)
	}
	dc := gg.NewContextForRGBA(p.image)

	grad := gg.NewRadialGradient(32, 46, 10, 32, 46, 80)
	grad.AddColorStop(0, color.RGBA{0, 255, 0, 255})
	grad.AddColorStop(1, color.RGBA{0, 0, 255, 255})

	dc.SetFillStyle(grad)
	dc.DrawRectangle(1, 1, 62, 94)
	dc.Fill()

	dc.SetColor(color.White)
	dc.DrawCircle(32, 48, 10)
	dc.Stroke()
	dc.DrawCircle(32, 48, 80)
	dc.Stroke()
	p.Update()
	<-time.After(time.Second)
}
