package display

import (
	"image"
	"image/draw"
	"sync"
)

type MirrorDisplay struct {
	displays []Display
}

func (m MirrorDisplay) Image() draw.Image {
	return m.displays[0].Image()
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

func (m MirrorDisplay) Update() {
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
