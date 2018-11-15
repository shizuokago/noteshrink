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

func TestRGB2HSV(t *testing.T) {
	//Black
	h, s, v := RGB2HSV(0, 0, 0)
	if !same(h, 0) || !same(s, 0) || !same(v, 0) {
		t.Errorf("Error:RGB2HSV Black value[%f][%f][%f]", h, s, v)
	}

	//White
	h, s, v = RGB2HSV(255, 255, 255)
	if !same(h, 0) || !same(s, 0) || !same(v, 1.0) {
		t.Errorf("Error:RGB2HSV White value.H=[%v],S[%f],V[%f]", h, s, v)
	}
	//Red
	h, s, v = RGB2HSV(255, 0, 0)
	if !same(h, 0) || !same(s, 1.0) || !same(v, 1.0) {
		t.Errorf("Error:RGB2HSV Red value.H=[%v],S[%f],V[%f]", h, s, v)
	}
	//Lime
	h, s, v = RGB2HSV(0, 255, 0)
	if !same(h, 0.333333) || !same(s, 1.0) || !same(v, 1.0) {
		t.Errorf("Error:RGB2HSV Lime value.H=[%v],S[%f],V[%f]", h, s, v)
	}
	//Blue
	h, s, v = RGB2HSV(0, 0, 255)
	if !same(h, 0.66667) || !same(s, 1.0) || !same(v, 1.0) {
		t.Errorf("Error:RGB2HSV Blue value.H=[%v],S[%f],V[%f]", h, s, v)
	}
	//Yellow
	h, s, v = RGB2HSV(255, 255, 0)
	if !same(h, 60.0/360.0) || !same(s, 1.0) || !same(v, 1.0) {
		t.Errorf("Error:RGB2HSV Yellow value.H=[%v],S[%f],V[%f]", h, s, v)
	}
	//Cyan
	h, s, v = RGB2HSV(0, 255, 255)
	if !same(h, 0.5) || !same(s, 1.0) || !same(v, 1.0) {
		t.Errorf("Error:RGB2HSV Cyan value.H=[%v],S[%f],V[%f]", h, s, v)
	}
	//Magenta
	h, s, v = RGB2HSV(255, 0, 255)
	if !same(h, 0.83333) || !same(s, 1.0) || !same(v, 1.0) {
		t.Errorf("Error:RGB2HSV Magenta value.H=[%v],S[%f],V[%f]", h, s, v)
	}
	//Silver
	h, s, v = RGB2HSV(192, 192, 192)
	//0.75???
	if !same(h, 0.0) || !same(s, 0.0) || !same(v, 0.752941) {
		t.Errorf("Error:RGB2HSV Sliver value.H=[%v],S[%f],V[%f]", h, s, v)
	}
	//Gray
	h, s, v = RGB2HSV(128, 128, 128)
	if !same(h, 0) || !same(s, 0) || !same(v, 0.501961) {
		t.Errorf("Error:RGB2HSV Gray value.H=[%v],S[%f],V[%f]", h, s, v)
	}
	//Maroon
	h, s, v = RGB2HSV(128, 0, 0)
	if !same(h, 0) || !same(s, 1.0) || !same(v, 0.501961) {
		t.Errorf("Error:RGB2HSV Maroon value.H=[%v],S[%f],V[%f]", h, s, v)
	}
	//Olive
	h, s, v = RGB2HSV(128, 128, 0)
	if !same(h, 60.0/360.0) || !same(s, 1.0) || !same(v, 0.501961) {
		t.Errorf("Error:RGB2HSV Olive value.H=[%v],S[%f],V[%f]", h, s, v)
	}
	//Green
	h, s, v = RGB2HSV(0, 128, 0)
	if !same(h, 120.0/360.0) || !same(s, 1.0) || !same(v, 0.501961) {
		t.Errorf("Error:RGB2HSV Green value.H=[%v],S[%f],V[%f]", h, s, v)
	}
	//Purple
	h, s, v = RGB2HSV(128, 0, 128)
	if !same(h, 300.0/360.0) || !same(s, 1.0) || !same(v, 0.501961) {
		t.Errorf("Error:RGB2HSV Purple value.H=[%v],S[%f],V[%f]", h, s, v)
	}
	//Teal
	h, s, v = RGB2HSV(0, 128, 128)
	if !same(h, 0.5) || !same(s, 1.0) || !same(v, 0.501961) {
		t.Errorf("Error:RGB2HSV Teal value.H=[%v],S[%f],V[%f]", h, s, v)
	}
	//Navy
	h, s, v = RGB2HSV(0, 0, 128)
	if !same(h, 240.0/360.0) || !same(s, 1.0) || !same(v, 0.501961) {
		t.Errorf("Error:RGB2HSV Navy value.H=[%v],S[%f],V[%f]", h, s, v)
	}
}

func TestHSV2RGB(t *testing.T) {
	//Black
	color := HSV2RGB(0, 0, 0)
	if color.R != 0 || color.G != 0 || color.B != 0 {
		t.Errorf("Error:HSV2RGB Black value[%v]", color)
	}

	//White
	color = HSV2RGB(0.0, 0.0, 1.0)
	if color.R != 255 || color.G != 255 || color.B != 255 {
		t.Errorf("Error:HSV2RGB White value[%v]", color)
	}
	//Red
	color = HSV2RGB(0, 1.0, 1.0)
	if color.R != 255 || color.G != 0 || color.B != 0 {
		t.Errorf("Error:HSV2RGB Red value[%v]", color)
	}
	//Lime
	color = HSV2RGB(0.33333, 1.0, 1.0)
	if color.R != 0 || color.G != 255 || color.B != 0 {
		t.Errorf("Error:HSV2RGB Lime value[%v]", color)
	}
	//Blue
	color = HSV2RGB(0.66667, 1.0, 1.0)
	if color.R != 0 || color.G != 0 || color.B != 255 {
		t.Errorf("Error:HSV2RGB Blue value[%v]", color)
	}
	//Yellow
	color = HSV2RGB(0.16667, 1.0, 1.0)
	if color.R != 255 || color.G != 255 || color.B != 0 {
		t.Errorf("Error:HSV2RGB Yellow value[%v]", color)
	}
	//Cyan
	color = HSV2RGB(0.5, 1.0, 1.0)
	if color.R != 0 || color.G != 255 || color.B != 255 {
		t.Errorf("Error:HSV2RGB Cyan value[%v]", color)
	}
	//Magenta
	color = HSV2RGB(0.83333, 1.0, 1.0)
	if color.R != 255 || color.G != 0 || color.B != 255 {
		t.Errorf("Error:HSV2RGB Magenta value[%v]", color)
	}
	//Silver
	color = HSV2RGB(0, 0, 0.75)
	//Not 192???
	if color.R != 191 || color.G != 191 || color.B != 191 {
		t.Errorf("Error:HSV2RGB Sliver value[%v]", color)
	}
	//Gray
	color = HSV2RGB(0, 0, 0.5)
	if color.R != 128 || color.G != 128 || color.B != 128 {
		t.Errorf("Error:HSV2RGB Gray value[%v]", color)
	}
	//Maroon
	color = HSV2RGB(0, 1.0, 0.5)
	if color.R != 128 || color.G != 0 || color.B != 0 {
		t.Errorf("Error:HSV2RGB Maroon value[%v]", color)
	}
	//Olive
	color = HSV2RGB(60.0/360.0, 1.0, 0.5)
	if color.R != 128 || color.G != 128 || color.B != 0 {
		t.Errorf("Error:HSV2RGB Olive value[%v]", color)
	}
	//Green
	color = HSV2RGB(120.0/360.0, 1.0, 0.5)
	if color.R != 0 || color.G != 128 || color.B != 0 {
		t.Errorf("Error:HSV2RGB Green value[%v]", color)
	}
	//Purple
	color = HSV2RGB(300.0/360.0, 1.0, 0.5)
	if color.R != 128 || color.G != 0 || color.B != 128 {
		t.Errorf("Error:HSV2RGB Purple value[%v]", color)
	}
	//Teal
	color = HSV2RGB(0.5, 1.0, 0.5)
	if color.R != 0 || color.G != 128 || color.B != 128 {
		t.Errorf("Error:HSV2RGB Teal value[%v]", color)
	}
	//Navy
	color = HSV2RGB(240.0/360.0, 1.0, 0.5)
	if color.R != 0 || color.G != 0 || color.B != 128 {
		t.Errorf("Error:HSV2RGB Navy value[%v]", color)
	}
}

func TestConvert(t *testing.T) {

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
		img.At(1000, 1000)
	}
}

func BenchmarkPack(b *testing.B) {
	p := NewPixelRGB(10, 20, 30)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Pack(p)
	}
}

func BenchmarkUnPack(b *testing.B) {
	p := NewPixelRGB(10, 20, 30)
	val := Pack(p)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		UnPack(val)
	}
}

func BenchmarkConvertColor(b *testing.B) {
	y, cb, cr := color.RGBToYCbCr(10, 20, 30)
	c := color.YCbCr{Y: y, Cb: cb, Cr: cr}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := convertColor(c)
		if err != nil {
			b.Errorf("convertColor() error [%v]", err)
			return
		}
	}
}

func BenchmarkRGB2HSV(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RGB2HSV(255, 255, 255)
	}
}

func BenchmarkHSV2RGB(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HSV2RGB(1, 1, 1)
	}
}

func BenchmarkFloatRGBA(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FloatRGBA(255, 255, 255)
	}
}

func BenchmarkUintRGBA(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		UIntRGBA(255, 255, 255)
	}
}

func same(a, b float64) bool {
	k := 0.00001
	if a >= b-k && a <= b+k {
		return true
	}
	return false
}
