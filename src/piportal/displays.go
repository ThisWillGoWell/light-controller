package piportal

import "github.com/thiswillgowell/light-controller/src/display"

var (
	TopLeftDisplay     display.Display
	TopRightDisplay    display.Display
	BottomLeftDisplay  display.Display
	BottomRightDisplay display.Display

	TopHalfDisplay    display.Display
	BottomHalfDisplay display.Display

	LeftHalfDisplay  display.Display
	RightHalfDisplay display.Display

	Fireplace display.Display
)

func init() {
	var err error
	topRightDisplay, err := NewMatrix("192.168.1.84:8080", TopLeft)
	if err != nil {
		panic(err)
	}
	TopRightDisplay = display.NewRotation(display.NewRotation(topRightDisplay, display.CounterClockwise), display.MirrorAcrossY)

	topLeftDisplay, err := NewMatrix("192.168.1.106:8080", BottomLeft)
	if err != nil {
		panic(err)
	}
	TopLeftDisplay = display.NewRotation(display.NewRotation(topLeftDisplay, display.Clockwise), display.MirrorAcrossY)

	bottomRightDisplay, err := NewMatrix("192.168.1.83:8080", TopRight)
	BottomRightDisplay = display.NewRotation(bottomRightDisplay, display.CounterClockwise)

	if err != nil {
		panic(err)
	}
	//	p2, err := NewMatrix("192.168.1.106:8080", TopRight)
	bottomLeftDisplay, err := NewMatrix("192.168.1.53:8080", TopLeft)
	if err != nil {
		panic(err)
	}
	BottomLeftDisplay = display.NewRotation(bottomLeftDisplay, display.Clockwise)

	BottomHalfDisplay = display.NewMultiDisplay(display.ArrangementHorizontal, BottomLeftDisplay, BottomRightDisplay)
	TopHalfDisplay = display.NewMultiDisplay(display.ArrangementHorizontal, TopLeftDisplay, TopRightDisplay)
	Fireplace = display.NewMultiDisplay(display.ArrangementVertical, BottomHalfDisplay, TopHalfDisplay)

}
