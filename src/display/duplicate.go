package display

import (
	"image"
	"image/color"
	"image/draw"
	"sync"
)

type DuplicateDisplay struct {
	displays []Display
}

func (m *DuplicateDisplay) ColorModel() color.Model {
	return m.displays[0].Image().ColorModel()
}

func (m *DuplicateDisplay) Bounds() image.Rectangle {
	return m.displays[0].Image().Bounds()

}

func (m *DuplicateDisplay) At(x, y int) color.Color {
	return m.displays[0].Image().At(x, y)
}

func (m *DuplicateDisplay) Set(x, y int, c color.Color) {
	for _, d := range m.displays {
		d.Image().Set(x, y, c)
	}
}

func (m *DuplicateDisplay) Image() draw.Image {
	return m
}

func (m DuplicateDisplay) UpdateImage(img image.Image) {
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

func (m *DuplicateDisplay) Update() {
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

func NewDuplicateDisplay(displays ...Display) *DuplicateDisplay {
	return &DuplicateDisplay{
		displays: displays,
	}
}
