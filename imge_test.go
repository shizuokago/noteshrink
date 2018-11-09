package noteshrink

import (
	"testing"

	"image/color"
)

func TestPack(t *testing.T) {

	c := color.RGBA{R: 255, G: 128, B: 10, A: 255}
	p := NewPixel(c)
	packed := Pack(p)
	ur, ug, ub := UnPack(packed)
	if c.R != ur || c.G != ug || c.B != ub {
		t.Errorf("Test Color:Pack [%v]!=[R:%d][G:%d][B:%d]", c, ur, ug, ub)
	}

	c = color.RGBA{R: 5, G: 5, B: 5, A: 255}
	p = NewPixel(c)
	packed = Pack(p)
	ur, ug, ub = UnPack(packed)
	if c.R != ur || c.G != ug || c.B != ub {
		t.Errorf("Test Color:Pack [%v]!=[R:%d][G:%d][B:%d]", c, ur, ug, ub)
	}
}

func TestColorConvert(t *testing.T) {

	t.Logf("RGB=[%d,%d,%d]", 1, 2, 255)
	h, s, v := RGB2HSV(1, 2, 255)
	if !same(h, 0.666010) || !same(s, 0.996078) || !same(v, 1.000000) {
		t.Errorf("Error:RGB2HSV Value[H:%f][S:%f][V:%f]", h, s, v)
	}

	newP := NewPixelHSV(h, s, v)
	if !same(h, newP.H) || !same(s, newP.S) || !same(v, newP.V) {
		t.Errorf("Error:NewPixelHSV value[%v]", newP)
	}

	t.Logf("Pixel=[%v]", newP)
	color := HSV2RGB(newP.H, newP.S, newP.V)
	if color.R != 1 || color.G != 2 || color.B != 255 {
		t.Errorf("Error:HSV2RGB value[%v]", color)
	}
}

func same(a, b float64) bool {
	k := 0.00001
	if a >= b-k && a <= b+k {
		return true
	}
	return false
}
