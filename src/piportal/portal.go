package piportal

import (
	"image"
	"image/draw"
)

type Matrix struct {
	image *image.RGBA
	Conn  *Connection
}

func (m *Matrix) Image() draw.Image {
	return m.image
}

func (m *Matrix) UpdateImage(src image.Image) {
	draw.Draw(m.image, m.image.Bounds(), src, image.Point{}, draw.Src)
	m.Update()
}

type PortalMode int

const (
	TopLeft PortalMode = iota
	TopRight
	BottomLeft
	BottomRight
)

func NewMatrix(address string, mode PortalMode) (*Matrix, error) {
	m := &Matrix{
		image: image.NewRGBA(image.Rect(0, 0, 64, 96)),
	}
	conn, err := NewTCPClient(address, mode)

	if err != nil {
		return nil, err
	}
	m.Conn = conn
	return m, nil
}

func (m *Matrix) Update() {
	m.Conn.WriteFrame(m.image)
}
