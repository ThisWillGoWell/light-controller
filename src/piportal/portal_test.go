package piportal

import (
	"github.com/thiswillgowell/light-controller/color"
	"testing"
)

func TestPortal(t *testing.T){

	p, err  := NewMatrix("192.168.1.63:8081")
	if err != nil {
		panic(err)
	}

	for i :=0;i<10000;i++{
		if err := p.ForEachAndUpdate(func(r, c int) color.Color {
			//return color.Deeppink
			return color.FromHsv(uint16(r*100 + c*100 + i*1000), 255, 30);
		}); err != nil {
			t.Fatal(err)
		}
	}
}








