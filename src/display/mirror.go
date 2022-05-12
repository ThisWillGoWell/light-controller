package display

import (
	"image"
	"image/color"
	"image/draw"
	"sync"
)

type MirrorDisplay struct {
	displays []Display
}

func (m *MirrorDisplay) ColorModel() color.Model {
	return m.displays[0].Image().ColorModel()
}

func (m *MirrorDisplay) Bounds() image.Rectangle {
	return m.displays[0].Image().Bounds()

}

func (m *MirrorDisplay) At(x, y int) color.Color {
	return m.displays[0].Image().At(x, y)
}

func (m *MirrorDisplay) Set(x, y int, c color.Color) {
	for _, d := range m.displays {
		d.Image().Set(x, y, c)
	}
}

func (m *MirrorDisplay) Image() draw.Image {
	return m
}

func (m MirrorDisplay) UpdateImage(img image.Image) {
	wg := sync.WaitGroup{}
	wg.Add(len(m.displays))
	for _, d := range m.displays {
		go func(d Display) {

			d.Update()
			wg.Done()
		}(d)
	}
	wg.Wait()
}

func (m *MirrorDisplay) Update() {
	wg := sync.WaitGroup{}
	wg.Add(len(m.displays))
	for _, d := range m.displays {
		go func(d Display) {
			d.Update()
			wg.Done()
		}(d)

	}
	wg.Wait()
}

func NewMirrorDisplay(displays ...Display) *MirrorDisplay {
	return &MirrorDisplay{
		displays: displays,
	}
}
