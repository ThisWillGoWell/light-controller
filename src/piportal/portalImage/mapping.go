package portalImage

// cords X1,Y1 X2,Y2 -> rotate

type PortalLayout int

const (
	TwoByThreeVertical PortalLayout = iota
)

func CreateMappedImage(layout PortalLayout) *Image {
	switch layout {
	case TwoByThreeVertical:
		return twoByThreeVerticalMappedPanel()
	}
	return nil
}

/*
twoByThreeVerticalMappedPanel create an image that maps the layout

|1a→|1b→|
|2a→|2b→|
|3a→|3b→|
to
|1b↑|2b↑|
|1a↑|2a↑|
|3b↑|3a↓|
*/
func twoByThreeVerticalMappedPanel() *Image {
	panelRows := 32
	panelColumns := 64

	mappedImage := NewMappedImage(panelRows*2, panelColumns*3, panelColumns*2, panelRows*3)

	// chain 1
	mappedImage.RegisterMapTransform(mappingTransform{
		startX: 0,
		startY: 0,
		endX:   panelColumns*2 - 1,
		endY:   panelRows - 1,
	}.define(func(x1, y1 int) (int, int) {
		return panelRows - y1 - 1, x1
	}))

	// chain 2
	mappedImage.RegisterMapTransform(mappingTransform{
		startX: 0,
		startY: panelRows,
		endX:   panelColumns*2 - 1,
		endY:   panelRows*2 - 1,
	}.define(func(x1, y1 int) (int, int) {
		return 3*panelRows - y1 - 2, x1
	}))

	// chain 3a
	mappedImage.RegisterMapTransform(mappingTransform{
		startX: panelColumns,
		startY: panelRows * 2,
		endX:   panelColumns*2 - 1,
		endY:   panelRows*3 - 1,
	}.define(func(x1, y1 int) (int, int) {
		return y1 - panelRows - 1, 4*panelColumns - x1 - 2
	}))

	// 3b
	mappedImage.RegisterMapTransform(mappingTransform{
		startX: 0,
		startY: panelRows * 2,
		endX:   panelColumns - 1,
		endY:   panelRows*3 - 1,
	}.define(func(x1, y1 int) (int, int) {
		return 3*panelRows - y1 - 1, 2*panelColumns + x1 - 1
	}))
	return mappedImage
}
