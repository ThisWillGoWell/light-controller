package text

import (
	"time"

	"github.com/hajimehoshi/bitmapfont/v2"
	"github.com/thiswillgowell/light-controller/src/display"
	"golang.org/x/image/font"
)

func ClockPattern(d display.Display) {
	f := bitmapfont.Face
	for {
		d.Clear()
		t := time.Now()
		tStr := t.Format("3:04")
		if t.Hour() > 12 {
			tStr += " AM"
		} else {
			tStr += " PM"
		}

		font.MeasureString(f, tStr)
		<-time.After(time.Second)
	}
}

func WriteFont(d display.Display) {

}
