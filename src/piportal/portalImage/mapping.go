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

	// chain 1 and 2
	mappedImage.RegisterMapTransform(mappingTransform{
		startX: 0,
		startY: 0,
		endY:   panelRows*2 - 1,
		endX:   panelColumns*2 - 1,
	}.rotateCCW().flipY())

	// chain 3a
	// rotate 180 degrees
	// move left, down one panel
	mappedImage.RegisterMapTransform(mappingTransform{
		startX: 0,
		startY: panelRows * 2,
		endX:   panelColumns,
		endY:   panelRows*3 - 1,
	}.rotateCW().translation(panelColumns, 0))

	mappedImage.RegisterMapTransform(mappingTransform{
		startX: panelColumns,
		endX:   panelColumns*2 - 1,
		startY: panelRows * 2,
		endY:   panelRows*3 - 1,
	}.rotateCCW().translation(-1*panelColumns, 0))

	return mappedImage
}
