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

	color := HSV2RGB(h, s, v)
	if color.R != 1 || color.G != 2 || color.B != 255 {
		t.Errorf("Error:HSV2RGB value[%v]", color)
	}
}

func BenchmarkImageAt(b *testing.B) {
	img, err := loadImage("sample/notesA1.jpg")
	if err != nil {
		b.Errorf("loadImage() Error[%v]", err)
		return
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		img.At(1000,1000)
	}
}

func BenchmarkPack(b *testing.B) {
	p := NewPixelRGB(10,20,30)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Pack(p)
	}
}

func BenchmarkUnPack(b *testing.B) {
	p := NewPixelRGB(10,20,30)
	val := Pack(p)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		UnPack(val)
	}
}

func BenchmarkConvertColor(b *testing.B) {
	y,cb,cr := color.RGBToYCbCr(10,20,30)
	c := color.YCbCr{Y:y,Cb:cb,Cr:cr}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := convertColor(c)
		if err != nil {
			b.Errorf("convertColor() error [%v]",err)
			return
		}
	}
}

func BenchmarkRGB2HSV(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RGB2HSV(255,255,255)
	}
}

func BenchmarkHSV2RGB(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HSV2RGB(1,1,1)
	}
}

func BenchmarkFloatRGBA(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FloatRGBA(255,255,255)
	}
}

func BenchmarkUintRGBA(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		UIntRGBA(255,255,255)
	}
}

func same(a, b float64) bool {
	k := 0.00001
	if a >= b-k && a <= b+k {
		return true
	}
	return false
}
