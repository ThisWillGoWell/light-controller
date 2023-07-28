package graphics

import (
	"math"
)

func MapFloatsToInts(values []float64, inMin, inMax float64, outMin, outMax, z int) []int {
	if z <= 0 {
		return nil
	}

	if len(values) == z {
		return mapInRange(values, inMin, inMax, outMin, outMax)
	}

	if len(values) > z {
		return mapToAverage(values, inMin, inMax, outMin, outMax, z)
	}

	return mapWithInterpolation(values, inMin, inMax, outMin, outMax, z)
}

func mapInRange(values []float64, inMin, inMax float64, outMin, outMax int) []int {
	outRange := outMax - outMin
	inRange := inMax - inMin

	result := make([]int, len(values))
	for i, v := range values {
		scaledValue := float64(outRange) * (v - inMin) / inRange
		result[i] = outMin + int(math.Round(scaledValue))
	}

	return result
}

func mapWithInterpolation(values []float64, inMin, inMax float64, outMin, outMax, z int) []int {
	result := make([]int, z)
	ratio := float64(len(values)-1) / float64(z-1)
	for i := 0; i < z; i++ {
		index := float64(i) * ratio
		lowerIndex := int(math.Floor(index))
		upperIndex := int(math.Ceil(index))
		lowerValue := values[lowerIndex]
		upperValue := values[upperIndex]

		inRange := inMax - inMin
		outRange := outMax - outMin

		interpolatedValue := lowerValue + (index-float64(lowerIndex))*(upperValue-lowerValue)
		scaledValue := float64(outRange) * (interpolatedValue - inMin) / inRange
		result[i] = outMin + int(math.Round(scaledValue))
	}

	return result
}

func mapToAverage(values []float64, inMin, inMax float64, outMin, outMax, z int) []int {
	if z <= 0 || len(values) == 0 {
		return []int{}
	}

	// Create slices to store the sum and count of values for each output index
	sumSlice := make([]float64, z)
	countSlice := make([]int, z)

	for i, v := range values {
		// Check if the value is within the input range
		if v >= inMin && v <= inMax {
			// Calculate the corresponding output index
			index := i * (z - 1) / (len(values) - 1) // Distribute input indices evenly over the output indices

			// Update the sum and count for the output index
			sumSlice[index] += v
			countSlice[index]++
		}
	}

	// Create the output list
	output := make([]int, z)
	for i := 0; i < z; i++ {
		// Calculate the average for each output index
		if countSlice[i] > 0 {
			average := sumSlice[i] / float64(countSlice[i])
			// Map the average to the output range
			output[i] = int(math.Round((average-inMin)/(inMax-inMin)*(float64(outMax)-float64(outMin)) + float64(outMin)))
		}
	}

	return output
}
