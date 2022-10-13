package daisy_test

import (
	"math"
	"testing"
)

func TestInitDaisy(t *testing.T) {
	//d, err := daisy.Init()
	//if err != nil {
	//	t.Fatal(err)
	//}
	////if err := music.StartDisplay(d.FFTChannel); err != nil {
	////	t.Fatal(err)
	////}
	//go func() {
	//	for floats := range d.FFTChannel {
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
