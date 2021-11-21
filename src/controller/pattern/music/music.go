package music

import (
	"math"
)


type MappingInput struct{
	InputSize int
	OutputSize int
	InputPercentages  []float64
	OutputPercentages []float64
}



type BinParams struct {
	NumBars int
	BarHeight int

	InputLength int

	LeftChannelFrequency *[]float32
	RightChannelFrequency *[]float32

	LeftChannelOutput *[]int
	RightChannelOutput *[]int
}

const maxFloatValue = float32(100.0)

// BuildIndexLUT  return a list of size param.InputSize that has each index of the output each index of the input goes to
func BuildIndexLUT(params MappingInput) []int{
	currentOutput := 0.0
	currentPercentageIndex := 0
	nextTransitionAt := int(math.Round(params.InputPercentages[0] * float64(params.InputSize)))
	outputPerInput :=  ( float64(params.OutputSize) * params.OutputPercentages[currentPercentageIndex]) / (float64(params.InputSize) * params.InputPercentages[currentPercentageIndex])

	lut := make([]int, params.InputSize)
	for i := range lut {
		lut[i] = -1
	}
	for i := 0; i< params.InputSize; i++ {
		if i == nextTransitionAt{
			// find when to next recalculate
			currentPercentageIndex += 1
			if currentPercentageIndex == len(params.InputPercentages) - 1 {
				nextTransitionAt = params.InputSize
			} else {
				nextTransitionAt = i + int(math.Round(params.InputPercentages[currentPercentageIndex] * float64(params.InputSize)))
			}

			outputPerInput =  ( float64(params.OutputSize) * params.OutputPercentages[currentPercentageIndex]) / (float64(params.InputSize) * params.InputPercentages[currentPercentageIndex])

		}

		if outputPerInput != 0 {
			lut[i] = int(currentOutput)
			currentOutput += outputPerInput
		}
	}
	return lut
}


type BinInput struct {
	Input []float32
	MaxInputValue float32
	MaxOutputValue int
	BinningLut []int
	NumOutputs int
}
func BinHeight( params BinInput) []int{
	output := make([]int, params.NumOutputs)
	outputPower := make([]float64, params.NumOutputs)
	inputsCountPerOutput := make([]int, params.NumOutputs)
	for i, inputVal := range params.Input {
		inputVal := math.Min(float64(params.MaxInputValue), float64(inputVal))
		index := params.BinningLut[i]
		if index == -1 {
			continue
		}
		outputPower[index] += math.Pow(inputVal, 2.0)
		inputsCountPerOutput[index] += 1
	}

	scalar := float64(float32(params.MaxOutputValue) / params.MaxInputValue)
	// calculate the RMS of each output power
	for i, outputPower := range outputPower {
		if inputsCountPerOutput[i] == 0 {
			continue
		}
		rms := math.Sqrt( 1 / float64(inputsCountPerOutput[i]) * outputPower)
		output[i] = int(rms * scalar)
	}
	return output
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