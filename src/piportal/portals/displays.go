package portals

import (
	"github.com/thiswillgowell/light-controller/src/display"
	"github.com/thiswillgowell/light-controller/src/piportal"
	"github.com/thiswillgowell/light-controller/src/piportal/portalImage"
)
import _ "embed"

var (
	LeftVertical  display.Display
	RightVertical display.Display
	BothVerticals display.Display
)

func init() {
	var err error
	leftVerticalRaw, err := piportal.NewMatrix("192.168.1.203:8080", portalImage.TwoByThreeVertical)
	if err != nil {
		panic(err)
	}
	LeftVertical = display.NewRotation(leftVerticalRaw, display.MirrorAcrossX)
	rightVerticalRaw, err := piportal.NewMatrix("192.168.1.200:8080", portalImage.TwoByThreeVertical)
	if err != nil {
		panic(err)
	}
	RightVertical = display.NewRotation(rightVerticalRaw, display.MirrorAcrossY)
	BothVerticals = display.NewDuplicateDisplay(LeftVertical, RightVertical)

}