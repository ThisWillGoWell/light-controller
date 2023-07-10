package portals

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
	TopRightDisplay = display.NewRotation(topRightDisplay, display.CounterClockwise)

	topLeftDisplay, err := NewMatrix("192.168.1.106:8080", BottomLeft)
	if err != nil {
		panic(err)
	}
	TopLeftDisplay = display.NewRotation(topLeftDisplay, display.Clockwise)

	bottomRightDisplay, err := NewMatrix("192.168.1.83:8080", TopRight)
	BottomRightDisplay = display.NewRotation(display.NewRotation(bottomRightDisplay, display.CounterClockwise), display.MirrorAcrossY)

	if err != nil {
		panic(err)
	}
	//	p2, err := NewMatrix("192.168.1.106:8080", TopRight)
	bottomLeftDisplay, err := NewMatrix("192.168.1.53:8080", TopLeft)
	if err != nil {
		panic(err)
	}
	BottomLeftDisplay = display.NewRotation(display.NewRotation(bottomLeftDisplay, display.Clockwise), display.MirrorAcrossY)

	BottomHalfDisplay = display.NewMultiDisplay(display.ArrangementHorizontal, BottomLeftDisplay, BottomRightDisplay)
	TopHalfDisplay = display.NewMultiDisplay(display.ArrangementHorizontal, TopLeftDisplay, TopRightDisplay)
	Fireplace = display.NewMultiDisplay(display.ArrangementVertical, TopHalfDisplay, BottomHalfDisplay)

}
