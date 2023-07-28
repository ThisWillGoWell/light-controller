package graphics

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapFloatsToInts_AllInInputRange(t *testing.T) {
	result := MapFloatsToInts([]float64{0.1, 0.5, 0.9}, 0, 1, 0, 10, 5)
	assert.Equal(t, []int{1, 0, 5, 0, 9}, result)
}
