package display

import (
	"image"
	"image/draw"
)

type MirrorDisplay struct {
	displays []Display
}

func (m MirrorDisplay) Image() draw.Image {
	return m.displays[0].Image()
}

func (m MirrorDisplay) UpdateImage(image image.Image) {
	for _, d := range m.displays {
		d.UpdateImage(image)
	}
}

func (m MirrorDisplay) Update() {
	for _, d := range m.displays {
		d.Update()
	}
}

func NewMirrorDisplay(displays ...Display) *MirrorDisplay {
	return &MirrorDisplay{
		displays: displays,
	}
}
