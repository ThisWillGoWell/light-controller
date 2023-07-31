package portals

import (
	"fmt"
	"github.com/thiswillgowell/light-controller/src/display"
	"github.com/thiswillgowell/light-controller/src/piportal"
	"github.com/thiswillgowell/light-controller/src/piportal/portalImage"
	"go.uber.org/zap"
)
import _ "embed"

var (
	LeftVertical  display.Display
	RightVertical display.Display
	BothVerticals display.Display
)
var logger = zap.S().With("package", "portals")

func init() {
	var err error
	fmt.Println("start")
	logger.Info("connecting to left panel")
	leftVerticalRaw, err := piportal.NewMatrix("192.168.1.84:8080", portalImage.TwoByThreeVertical)
	if err != nil {
		panic(err)
	}
	fmt.Println("end1")
	LeftVertical = display.NewRotation(display.NewRotation(leftVerticalRaw, display.MirrorAcrossY), display.MirrorAcrossX)
	logger.Info("connected")
	logger.Info("connecting to right panel")
	RightVertical, err = piportal.NewMatrix("192.168.1.106:8080", portalImage.TwoByThreeVertical)
	if err != nil {
		panic(err)
	}
	logger.Info("connected")
	BothVerticals = display.NewDuplicateDisplay(LeftVertical, RightVertical)
	fmt.Println("end")
}
