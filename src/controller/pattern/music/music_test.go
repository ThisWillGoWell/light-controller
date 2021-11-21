package music

import "testing"
import   "github.com/stretchr/testify/assert"

func TestBuildIndexLUT(t *testing.T){
	t.Run("simple input 1-1", func(t *testing.T) {
		input := MappingInput{
			InputSize: 10,
			OutputSize: 10,
			InputPercentages: []float64{1},
			OutputPercentages:  []float64{1},
		}
		output := []int{0,1,2,3,4,5,6,7,8,9}

		lut := BuildIndexLUT(input)
		assert.Equal(t, output, lut)
	})

	t.Run("simple input 2-1", func(t *testing.T) {
		input := MappingInput{
			InputSize: 20,
			OutputSize: 10,
			InputPercentages: []float64{1},
			OutputPercentages:  []float64{1},
		}
		output := []int{0,0,1,1,2,2,3,3,4,4,5,5,6,6,7,7,8,8,9,9}

		lut := BuildIndexLUT(input)
		assert.Equal(t, output, lut)
	})

	t.Run("0.5->1", func(t *testing.T) {
		input := MappingInput{
			InputSize: 10,
			OutputSize: 10,
			InputPercentages: []float64{0.5, 0.5},
			OutputPercentages:  []float64{1, 0},
		}
		output := []int{0,2,4,6,8, -1, -1, -1, -1, -1}

		lut := BuildIndexLUT(input)
		assert.Equal(t, output, lut)
	})

	t.Run("complexOne", func(t *testing.T) {
		input := MappingInput{
			InputSize: 10,
			OutputSize: 10,
			InputPercentages: []float64{0.30, 0.60},
			OutputPercentages:  []float64{0.5, 0.5},
		}
		output := []int{0,2,4,6,8, -1, -1, -1, -1, -1}

		lut := BuildIndexLUT(input)
		assert.Equal(t, output, lut)
	})
}
