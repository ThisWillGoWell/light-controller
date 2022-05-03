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
	topRightDisplay, err := NewMatrix("192.168.1.84:8080", Right)
	if err != nil {
		panic(err)
	}
	TopRightDisplay = display.NewRotation(topRightDisplay, display.MirrorAcrossX)

	TopLeftDisplay, err = NewMatrix("192.168.1.106:8080", Left)
	if err != nil {
		panic(err)
	}

	BottomRightDisplay, err = NewMatrix("192.168.1.83:8080", Left)
	if err != nil {
		panic(err)
	}
	//	p2, err := NewMatrix("192.168.1.106:8080", Left)
	bottomLeftDisplay, err := NewMatrix("192.168.1.53:8080", Right)
	if err != nil {
		panic(err)
	}
	BottomLeftDisplay = display.NewRotation(bottomLeftDisplay, display.MirrorAcrossY)

	BottomHalfDisplay = display.NewMultiDisplay(display.ArrangementVertical, BottomLeftDisplay, BottomRightDisplay)

	Fireplace = display.NewMultiDisplay(display.ArrangementVertical, BottomHalfDisplay, TopHalfDisplay)

}
