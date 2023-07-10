package piportal

import (
	"github.com/thiswillgowell/light-controller/src/piportal/portalImage"
	"image"
	"image/draw"
)

type Matrix struct {
	image *portalImage.Image
	Conn  *Connection
}

func (m *Matrix) Image() draw.Image {
	return m.image
}

func (m *Matrix) UpdateImage(src image.Image) {
	draw.Draw(m.image, m.image.Bounds(), src, image.Point{}, draw.Src)
	m.Update()
}

func NewMatrix(address string, layout portalImage.PortalLayout) (*Matrix, error) {
	m := &Matrix{
		image: portalImage.CreateMappedImage(layout),
	}
	conn, err := NewTCPClient(address)

	if err != nil {
		return nil, err
	}
	m.Conn = conn
	return m, nil
}

func (m *Matrix) Update() {
	m.Conn.WriteFrame(m.image.MappedImage())
}
