package specturm

import "math"

// BinInput given an input of FrequencyValue, convert
// to set of integers
type BinInput struct {
	// what is the max expected value of the input
	MaxInputValue FrequencyValue
	// max integer to bin the
	MaxOutputValue int
	// A set of binning
	BinningLut []int

	NumOutputs int
	// if you are binning to more outputs then inputs, this will interpolate
	// the values between values
	Interpolate bool
}

func BinHeight(input []FrequencyValue, params BinInput, maxInputValue FrequencyValue) []int {
	output := make([]int, params.NumOutputs)
	outputPower := make([]float64, params.NumOutputs)
	inputsCountPerOutput := make([]int, params.NumOutputs)

	hasHardValue := make([]bool, params.NumOutputs)

	for i, inputVal := range input {
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
	scalar := float64(float32(params.MaxOutputValue) / float32(params.MaxInputValue))

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

// MappingInput create the BinningLut
// allow for non-linar mapping of input frequencies to
// bins
type MappingInput struct {
	// number of inputs
	InputSize int
	// number of outputs to map too
	OutputSize int
	// given a list of floats that must add to 1
	// to say what input percentages map to the corresponding
	// output percentages
	InputPercentages  []float64
	OutputPercentages []float64

	// map Input[InputSize-1] to Output[0]
	Reversed bool
}

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
