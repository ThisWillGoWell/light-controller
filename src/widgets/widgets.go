package widgets

import (
	"github.com/thiswillgowell/light-controller/src/display"
	"image"
	"image/draw"
	"time"
)

type Widget interface {
	Draw(img draw.Image, p image.Point)
}

type WidgetManager struct {
	disp       display.Display
	widgets    map[string]Widget
	positions  map[string]image.Point
	updateRate time.Duration
}

func (wm *WidgetManager) Update() {
	for name, widget := range wm.widgets {
		widget.Draw(wm.disp.Image(), wm.positions[name])
	}
}
