package music

import (
	"image/color"
	"math"

	"github.com/thiswillgowell/light-controller/src/graphics"
	"golang.org/x/image/colornames"

	"github.com/fogleman/gg"

	"github.com/thiswillgowell/light-controller/src/daisy/daisy"
	"github.com/thiswillgowell/light-controller/src/display"
)

type MappingInput struct {
	InputSize         int
	OutputSize        int
	InputPercentages  []float64
	OutputPercentages []float64
	Reversed          bool
}

type BinParams struct {
	NumBars   int
	BarHeight int

	InputLength int

	LeftChannelFrequency  *[]float32
	RightChannelFrequency *[]float32

	LeftChannelOutput  *[]int
	RightChannelOutput *[]int
}

const maxFloatValue = float32(100.0)

// BuildIndexLUT  return a list of size param.InputSize that has each index of the output each index of the input goes to
func BuildIndexLUT(params MappingInput) []int {
	currentOutput := 0.0
	currentPercentageIndex := 0
	nextTransitionAt := int(math.Round(params.InputPercentages[0] * float64(params.InputSize)))
	outputPerInput := (float64(params.OutputSize) * params.OutputPercentages[currentPercentageIndex]) / (float64(params.InputSize) * params.InputPercentages[currentPercentageIndex])

	lut := make([]int, params.InputSize)
	for i := range lut {
		lut[i] = -1
	}
	for i := 0; i < params.InputSize; i++ {
		if i == nextTransitionAt {
			// find when to next recalculate
			currentPercentageIndex += 1
			if currentPercentageIndex == len(params.InputPercentages)-1 {
				nextTransitionAt = params.InputSize
			} else {
				nextTransitionAt = i + int(math.Round(params.InputPercentages[currentPercentageIndex]*float64(params.InputSize)))
			}

			outputPerInput = (float64(params.OutputSize) * params.OutputPercentages[currentPercentageIndex]) / (float64(params.InputSize) * params.InputPercentages[currentPercentageIndex])

		}

		if outputPerInput != 0 {
			if params.Reversed {
				lut[i] = params.OutputSize - int(currentOutput) - 1
			} else {
				lut[i] = int(currentOutput)
			}

			currentOutput += outputPerInput
		}
	}
	return lut
}

type BinInput struct {
	Input          []float32
	MaxInputValue  float32
	MaxOutputValue int
	BinningLut     []int
	NumOutputs     int

	Interpolate bool
}

func BinHeight(params BinInput) []int {
	output := make([]int, params.NumOutputs)
	outputPower := make([]float64, params.NumOutputs)
	inputsCountPerOutput := make([]int, params.NumOutputs)

	hasHardValue := make([]bool, params.NumOutputs)

	for i, inputVal := range params.Input {
		inputVal := math.Min(float64(params.MaxInputValue), float64(inputVal))
		index := params.BinningLut[i]
		if index == -1 {
			continue
		}
		// used for rms vals
		//outputPower[index] += math.Pow(inputVal, 2.0)
		outputPower[index] = math.Max(outputPower[index], inputVal)
		hasHardValue[index] = true
		inputsCountPerOutput[index] += 1
	}
	scalar := float64(float32(params.MaxOutputValue) / params.MaxInputValue)
	// calculate the RMS of each output power
	for i, power := range outputPower {
		if inputsCountPerOutput[i] == 0 {
			continue
		}
		//rms := math.Sqrt(1 / float64(inputsCountPerOutput[i]) * outputPower)
		output[i] = int(power * scalar)
	}

	if params.Interpolate {
		lastKnownValue := output[0]
		backTrackCount := 0
		for i, knownValue := range hasHardValue {
			if !knownValue {
				backTrackCount += 1
				continue
			}
			if backTrackCount == 0 {
				continue
			}
			currentValue := output[i]
			// linear interpolation
			scalar := float64((currentValue - lastKnownValue) / (backTrackCount + 1))
			for backTrackCount != 0 {
				output[i-backTrackCount] = int(float64(currentValue) - (scalar * float64(backTrackCount)))
				backTrackCount -= 1
			}
			lastKnownValue = currentValue
		}
	}
	return output
}

type BinController struct {
	maxPower int
}

func (bc *BinController) SmartBin(params BinInput) []int {
	//output := make([]int, params.NumOutputs)
	//outputPower := make([]float64, params.NumOutputs)
	//inputsCountPerOutput := make([]int, params.NumOutputs)
	//
	//hasHardValue := make([]bool, params.NumOutputs)
	//
	//for i, inputVal := range params.Input {
	//	if inputVal > bc.maxPower {
	//		bc.maxPower =
	//	}
	//
	//	inputVal := math.Min(float64(params.MaxInputValue), float64(inputVal))
	//	index := params.BinningLut[i]
	//	if index == -1 {
	//		continue
	//	}
	//	// used for rms vals
	//	//outputPower[index] += math.Pow(inputVal, 2.0)
	//	if outputPower[]
	//	outputPower[index] = math.Max(outputPower[index], maxPower)
	//	hasHardValue[index] = true
	//	inputsCountPerOutput[index] += 1
	//}
	//scalar := float64(float32(params.MaxOutputValue) / params.MaxInputValue)
	//// calculate the RMS of each output power
	//for i, power := range outputPower {
	//	if inputsCountPerOutput[i] == 0 {
	//		continue
	//	}
	//	//rms := math.Sqrt(1 / float64(inputsCountPerOutput[i]) * outputPower)
	//	output[i] = int(power * scalar)
	//}
	//
	//if params.Interpolate {
	//	lastKnownValue := output[0]
	//	backTrackCount := 0
	//	for i, knownValue := range hasHardValue {
	//		if !knownValue {
	//			backTrackCount += 1
	//			continue
	//		}
	//		if backTrackCount == 0 {
	//			continue
	//		}
	//		currentValue := output[i]
	//		// linear interpolation
	//		scalar := float64((currentValue - lastKnownValue) / (backTrackCount + 1))
	//		for backTrackCount != 0 {
	//			output[i-backTrackCount] = int(float64(currentValue) - (scalar * float64(backTrackCount)))
	//			backTrackCount -= 1
	//		}
	//		lastKnownValue = currentValue
	//	}
	//}
	//return output
	return nil
}

//// floatsPerBin Given a index of the frequency, return
//// give half the floats 2/3's of the bars
//func floatsPerBar(index, totalFrequencies, numBars int) int {
//	// 1/3 of the bars
//	floatsPerBin1 := 1 / (float32(totalFrequencies) / (float32(numBars) * 2.0 / 3.0))
//	if index > totalFrequencies / 2 {
//		// inverse of frequency per bar,
//		return
//	}else {
//		return 1/2
//	}
//}
//
//func binFrequency(frequencies *[]float32, output *[]int, barHeight int  ) {
//	floatsPerBin :=
//}
//
//func BinFrequencies(ctx context.Context, params BinParams ){
//	barHeight := make([]int, params.BarHeight)
//	floatsPerBin1 := 1 / ((float32(params.InputLength)) / (float32(len(barHeight)) * 2.0 / 3.0)) // # half the floats get 2/3s of the bins
//	floatsPerBin2 := floatsPerBin1 / 2.0
//	bin := 0
//
//	for i, val := range params.TwoChanelFrequency {
//		if i >= len(params.TwoChanelFrequency)/2 {
//			if i < len(params.TwoChanelFrequency)*3/4 { // first half, second channel
//				bin = len(barHeight)/2 + int(floatsPerBin1*float32(i-len(params.TwoChanelFrequency)/2))
//			} else { // second half, second channel
//				bin = len(barHeight)*5/6 + int(floatsPerBin2*(float32(i)-float32(len(params.TwoChanelFrequency))/4*3))
//			}
//		} else {
//			if i < len(params.TwoChanelFrequency)/4 {
//				bin = int(floatsPerBin1 * float32(i))
//			} else {
//				bin = len(barHeight)*1/3 + int(floatsPerBin2*(float32(i)-float32(len(params.TwoChanelFrequency))/4))
//			}
//		}
//
//		if bin >= len(barHeight) {
//			bin = len(barHeight) - 1
//		}
//		//bin := int(float32(i) / float32(len(floats)) * float32(matrix.Cols*2))
//		// ,:,l;val = val * float32(math.Pow(2, float64(i)))
//		if val > maxFloatValue {
//			val = maxFloatValue
//		}
//		//scale := 1 - 2/(math.Exp(float64(2*val/maxFloatValue))+1)
//		//if scale < 0 {
//		//	scale = 0
//		//}
//		val = val / maxFloatValue
//		size := int(val * float32(m.Rows/2))
//		if val < 0 {
//			size = 0
//		}
//		if size > barHeight[bin] {
//			barHeight[bin] = size
//		}
//	}
//	return barHeight
//
//}
//
//func CenterVUBar(daisyDevice *daisy.Daisy, display display.Display) {
//	colors := color_2.LinerGradient(color_2.Darkred, color_2.Purple, display.Cols(), false)
//	lut := BuildIndexLUT(MappingInput{
//		InputSize:         daisy.NumFrequencies,
//		OutputSize:        display.Cols(),
//		InputPercentages:  []float64{0.4, 0.3, 0.3},
//		OutputPercentages: []float64{0.5, 0.4, 0.1},
//		Reversed:          true,
//	})
//
//	for channel := range daisyDevice.FFTChannel {
//		controller.ForEach(display, controller.DarkenDisplay(.3))
//		leftBars := BinHeight(BinInput{
//			Input:          channel[0],
//			MaxInputValue:  200.0,
//			MaxOutputValue: display.Rows() / 2,
//			BinningLut:     lut,
//			NumOutputs:     display.Cols(),
//		})
//		// left is top half of the display
//		startRow := display.Rows() / 2
//		for col, height := range leftBars {
//			for i := 0; i < height; i++ {
//				display.SetPixel(startRow+i, col, colors[col])
//			}
//		}
//
//		rightBars := BinHeight(BinInput{
//			Input:          channel[1],
//			MaxInputValue:  200.0,
//			MaxOutputValue: display.Rows() / 2,
//			BinningLut:     lut,
//			NumOutputs:     display.Cols(),
//		})
//
//		// right is top half of the display, but looking down
//		startRow = display.Rows()/2 - 1
//		for col, height := range rightBars {
//			for i := 0; i < height; i++ {
//				display.SetPixel(startRow-i, col, colors[col])
//			}
//		}
//		display.Send()
//	}
//}
//
//func CenterHollowVUBar(daisyDevice *daisy.Daisy, display display.Display) {
//	colors := color_2.LinerGradient(color_2.Darkred, color_2.Purple, display.Cols(), false)
//	lut := BuildIndexLUT(MappingInput{
//		InputSize:         daisy.NumFrequencies,
//		OutputSize:        display.Cols(),
//		InputPercentages:  []float64{0.4, 0.3, 0.3},
//		OutputPercentages: []float64{0.5, 0.4, 0.1},
//		Reversed:          true,
//	})
//	image :=
//	for channel := range daisyDevice.FFTChannel {
//		controller.ForEach(display, controller.DarkenDisplay(.2))
//		//display.Clear()
//		leftBars := BinHeight(BinInput{
//			Input:          channel[0],
//			MaxInputValue:  200.0,
//			MaxOutputValue: display.Rows() / 2,
//			BinningLut:     lut,
//			NumOutputs:     display.Cols(),
//		})
//		// left is top half of the display
//		startRow := display.Rows() / 2
//		for col, height := range leftBars {
//			for i := 0; i < height; i++ {
//				c := colors[col].DarkenPercentage(float64(i+2) / float64(height))
//				display.SetPixel(startRow+i, col, c)
//			}
//		}
//
//		rightBars := BinHeight(BinInput{
//			Input:          channel[1],
//			MaxInputValue:  200.0,
//			MaxOutputValue: display.Rows() / 2,
//			BinningLut:     lut,
//			NumOutputs:     display.Cols(),
//		})
//
//		// right is top half of the display, but looking down
//		startRow = display.Rows()/2 - 1
//		for col, height := range rightBars {
//			for i := 0; i < height; i++ {
//				c := colors[col].DarkenPercentage(float64(i+2) / float64(height))
//				display.SetPixel(startRow-i, col, c)
//			}
//		}
//		display.Send()
//	}
//}

func CenterHollowVUBarDouble(daisyDevice interface{ NextFFTValues() [][]float32 }, d display.Display, barWidth int) {
	targetImage := d.Image()
	numBars := targetImage.Bounds().Size().X / barWidth
	colors := graphics.LinerGradient(colornames.Purple, colornames.Darkorange, numBars)
	lut := BuildIndexLUT(MappingInput{
		InputSize:         daisy.NumFrequencies,
		OutputSize:        numBars,
		InputPercentages:  []float64{0.4, 0.3, 0.3},
		OutputPercentages: []float64{0.5, 0.4, 0.1},
		Reversed:          false,
	})

	gg.NewContext(targetImage.Bounds().Size().X, targetImage.Bounds().Size().Y)
	dc := gg.NewContextForImage(targetImage)
	for {
		channel := daisyDevice.NextFFTValues()
		if channel == nil {
			break
		}
		dc.SetColor(color.Transparent)
		dc.Clear()
		leftBars := removeDeadZones(BinHeight(BinInput{
			Input:          channel[0],
			MaxInputValue:  200.0,
			MaxOutputValue: targetImage.Bounds().Size().Y / 2 * 5 / 6,
			BinningLut:     lut,
			NumOutputs:     numBars,
			Interpolate:    true,
		}))
		// left is top half of the display
		startRow := targetImage.Bounds().Size().Y / 2

		for col, barPower := range leftBars {
			c := colors[col]
			// 0 - 30%
			var barColors []color.Color
			if barPower < 3 {
				c = graphics.Darken(c, 1-float64(barPower)*.1)
				barColors = []color.Color{c}
			} else {
				// 30% -> 100%
				barColors = graphics.LinerGradient(graphics.Darken(c, .7), c, barPower-2)
			}

			for i, c := range barColors {
				dc.SetColor(c)
				for j := 0; j < barWidth; j++ {
					dc.SetPixel(col*barWidth+j, startRow+i)
				}
			}
		}

		rightBars := removeDeadZones(BinHeight(BinInput{
			Input:          channel[1],
			MaxInputValue:  200.0,
			MaxOutputValue: targetImage.Bounds().Size().Y / 2 * 5 / 6,
			BinningLut:     lut,
			NumOutputs:     numBars,
			Interpolate:    true,
		}))

		// right is top half of the display, but looking down
		for col, barPower := range rightBars {
			c := colors[col]
			// 0 - 30%
			var barColors []color.Color
			if barPower < 3 {
				c = graphics.Darken(c, 1-float64(barPower)*.1)
				barColors = []color.Color{c}
			} else {
				// 30% -> 100%
				barColors = graphics.LinerGradient(graphics.Darken(c, .7), c, barPower-2)
			}
			for i, c := range barColors {
				dc.SetColor(c)
				for j := 0; j < barWidth; j++ {
					dc.SetPixel(col*barWidth+j, startRow-i-1)
				}
			}
		}
		display.DrawAndUpdate(d, dc.Image())
		d.Update()
	}
}

//
//func FallingVuMeter(daisyDevice *daisy.Daisy, display display.Display) {
//	colors := color_2.LinerGradient(color_2.Darkred, color_2.ForestGreen, display.Rows(), true)
//	lut := BuildIndexLUT(MappingInput{
//		InputSize:         daisy.NumFrequencies,
//		OutputSize:        display.Cols() / 2,
//		InputPercentages:  []float64{0.4, 0.3, 0.3},
//		OutputPercentages: []float64{0.5, 0.4, 0.1},
//		Reversed:          false,
//	})
//
//	for channel := range daisyDevice.FFTChannel {
//		controller.ForEach(display, controller.DarkenDisplay(.3))
//		leftBars := BinHeight(BinInput{
//			Input:          channel[0],
//			MaxInputValue:  200.0,
//			MaxOutputValue: display.Rows(),
//			BinningLut:     lut,
//			NumOutputs:     display.Cols() / 2,
//		})
//		// left is top half of the display
//		startRow := display.Rows() / 2
//		for col, height := range leftBars {
//			c := colors[col]
//			if height < 4 {
//				c = c.DarkenPercentage(float64(5-height) * .1)
//			}
//
//			for i := 0; i < height; i++ {
//				c := c.DarkenPercentage(float64(i) / float64(height))
//				display.SetPixel(startRow+i, col*2, c)
//				display.SetPixel(startRow+i, col*2+1, c)
//			}
//		}
//
//		rightBars := BinHeight(BinInput{
//			Input:          channel[1],
//			MaxInputValue:  200.0,
//			MaxOutputValue: display.Rows() / 2,
//			BinningLut:     lut,
//			NumOutputs:     display.Cols() / 2,
//		})
//
//		// right is top half of the display, but looking down
//		startRow = display.Rows()/2 - 1
//		for col, height := range rightBars {
//			c := colors[col]
//			if height < 4 {
//				c = c.DarkenPercentage(float64(5-height) * .1)
//			}
//			for i := 0; i < height; i++ {
//				c := c.DarkenPercentage(float64(i) / float64(height))
//				display.SetPixel(startRow-i, col*2, c)
//				display.SetPixel(startRow-i, col*2+1, c)
//			}
//		}
//		display.Send()
//	}
//}

func removeDeadZones(input []int) []int {
	for i, ele := range input {
		if i == 0 || i == len(input)-1 {
			continue
		}
		// leave alone if bigger than one
		if ele > 1 {
			continue
		}
		// set to zero if previous value < 1 and next value is < 1

		if input[i-1] <= 1 && input[i+1] <= 1 {
			input[i] = 0
		}
	}
	return input
}
