package daisy

import "math/rand"

type Mock struct {
	previousValues [][]float32
}

// NextFFTValues will take the current values, and go plus or minus a random precentage
func (m *Mock) NextFFTValues() [][]float32 {
	if m.previousValues == nil {
		m.previousValues = make([][]float32, 2)
		m.previousValues[0] = make([]float32, NumFrequencies)
		m.previousValues[1] = make([]float32, NumFrequencies)
	}
	for i, val := range m.previousValues[0] {
		// scale +/- 5%
		scaleAmount := float32(rand.Intn(10) - 5)
		val += scaleAmount
		if val < 0 {
			val = 100
		}
		if val > 200 {
			val = 100
		}

		m.previousValues[0][i] = val
	}

	for i, val := range m.previousValues[1] {
		// scale +/- 5%
		scaleAmount := float32(rand.Intn(10) - 5)
		val += scaleAmount
		if val < 0 {
			val = 100
		}
		if val > 200 {
			val = 100
		}

		m.previousValues[1][i] = val
	}

	return m.previousValues
}
