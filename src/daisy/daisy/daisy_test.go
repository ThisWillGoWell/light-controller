package daisy_test

import (
	"math"
	"testing"
)

func TestInitDaisy(t *testing.T) {
	//d, err := daisy.InitDaisy()
	//if err != nil {
	//	t.Fatal(err)
	//}
	//if err := music.StartDisplay(d.FFTChannel); err != nil {
	//	t.Fatal(err)
	//}
	//go func() {
	//	for floats := range d.FFTChannel {
	//
	//		maxFrequencyLeft := 0
	//		maxFrequencyRight := 0
	//		maxFloatLeft := float32(0)
	//		maxFloatRight := float32(0)
	//
	//		for i, float := range floats {
	//
	//			if i > len(floats)/2 {
	//				if float > maxFloatRight {
	//					maxFloatRight = float
	//					maxFrequencyRight = i
	//				}
	//			} else {
	//				if float > maxFloatLeft {
	//					maxFloatLeft = float
	//					maxFrequencyLeft = i
	//				}
	//			}
	//		}
	//
	//		fmt.Printf("%d\t%4f\t%d\t%4f\n", maxFrequencyLeft, maxFloatLeft, maxFrequencyRight, maxFloatRight)
	//		<-time.After(time.Millisecond * 100)
	//	}
	//}()
	//for true {
	//
	//}
}

func PanioFreq(i int) float32 {
	return float32(math.Round(math.Pow(math.Pow(2, 1.0/12.0), float64(i-49)) * 440))
}
