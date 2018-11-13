package noteshrink

import (
	"testing"
	"math/rand"
)

func TestNewRGB(t *testing.T) {
	p := NewPixelRGB(50, 100, 201)
	if p.R != 50 || p.G != 100 || p.B != 201 {
		t.Errorf("NewPixelRGB[%v]", p)
	}
	//HSV
	//p.value
}

func TestRGB(t *testing.T) {

	p := NewPixelRGB(50, 100, 201)
	if p.R != 50 || p.G != 100 || p.B != 201 {
		t.Errorf("NewPixelRGB[%v]", p)
	}
	r, g, b := p.RGB()
	if p.R != r || p.G != g || p.B != b {
		t.Errorf("[%v] RGB()[%d][%d][%d]", p, r, g, b)
	}

}

func TestShift(t *testing.T) {

	p := NewPixelRGB(50, 100, 201)
	shift := p.Shift(2)

	if shift.R != 48 {
		t.Errorf("Shift R not 0[%d]", shift.R)
	}
	if shift.G != 100 {
		t.Errorf("Shift G not 0[%d]", shift.G)
	}
	if shift.B != 200 {
		t.Errorf("Shift B not 0[%d]", shift.B)
	}

	shift = p.Shift(3)
	if shift.R != 48 {
		t.Errorf("Shift R not 0[%d]", shift.R)
	}
	if shift.G != 96 {
		t.Errorf("Shift G not 0[%d]", shift.G)
	}
	if shift.B != 200 {
		t.Errorf("Shift B not 0[%d]", shift.B)
	}
}

func TestDistanceRGB(t *testing.T) {
	p1 := NewPixelRGB(100, 100, 100)
	p2 := NewPixelRGB(50, 50, 50)
	val := p1.DistanceRGB(p2)
	if val != 7500 {
		t.Errorf("DistanceRGB Error 7500!=[%f]", val)
	}

	p2 = NewPixelRGB(90, 90, 100)
	val = p1.DistanceRGB(p2)
	if val != 200 {
		t.Errorf("DistanceRGB Error 200!=[%f]", val)
	}

	p2 = NewPixelRGB(100, 90, 90)
	val = p1.DistanceRGB(p2)
	if val != 200 {
		t.Errorf("DistanceRGB Error 200!=[%f]", val)
	}

	p2 = NewPixelRGB(90, 100, 90)
	val = p1.DistanceRGB(p2)
	if val != 200 {
		t.Errorf("DistanceRGB Error 200!=[%f]", val)
	}
}

func TestDistanceHSV(t *testing.T) {
}

func TestMost(t *testing.T) {
	pix := make(Pixels, 0)
	pix = append(pix, NewPixelRGB(100, 100, 100))
	pix = append(pix, NewPixelRGB(100, 100, 100))
	pix = append(pix, NewPixelRGB(100, 100, 100))

	pix = append(pix, NewPixelRGB(50, 100, 100))
	pix = append(pix, NewPixelRGB(50, 100, 100))

	p := pix.Most()
	if p.R != 100 || p.G != 100 || p.B != 100 {
		t.Errorf("Most() Error")
	}
}

func TestQuantize(t *testing.T) {
	pix := make(Pixels, 0)
	pix = append(pix, NewPixelRGB(25, 50, 100))
	pix = append(pix, NewPixelRGB(150, 175, 200))

	newPix, err := pix.Quantize(5)
	if err != nil {
		t.Errorf("Quantize Error:[%v]", err)
	}

	if len(newPix) != 2 {
		t.Errorf("Quantize Length Error:not 2[%d]", len(newPix))
	}

	if newPix[0].R != 0 || newPix[0].G != 32 || newPix[0].B != 96 {
		t.Errorf("Quantize Error index 0[%v]", newPix[0])
	}

	if newPix[1].R != 128 || newPix[1].G != 160 || newPix[1].B != 192 {
		t.Errorf("Quantize Error index 1[%v]", newPix[1])
	}
}

func TestAverage(t *testing.T) {
	pix := make(Pixels, 0)
	pix = append(pix, NewPixelRGB(25, 50, 100))
	pix = append(pix, NewPixelRGB(150, 175, 200))

	newPix, err := pix.Average()

	if err != nil {
		t.Errorf("Average return error[%v]", err)
	} else {
		if newPix.R != 88 || newPix.G != 113 || newPix.B != 150 {
			t.Errorf("Quantize Error index 0[%v]", newPix)
		}
	}

	//nil error
}

func TestSort(t *testing.T) {
}

func BenchmarkNewPixel(b *testing.B) {
	img, err := loadImage("sample/notesA1.jpg")
	if err != nil {
		b.Errorf("loadImage() Error[%v]", err)
		return
	}
	c := img.At(1000,1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewPixel(c)
	}
}

func BenchmarkDistanceHSV(b *testing.B) {
	op := NewPixelRGB(100,100,100)
	sp := NewPixelRGB(200,200,200)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		op.DistanceHSV(sp)
	}
}

func BenchmarkDistanceRGB(b *testing.B) {
	op := NewPixelRGB(100,100,100)
	sp := NewPixelRGB(200,200,200)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		op.DistanceRGB(sp)
	}
}

func BenchmarkShift(b *testing.B) {
	p := NewPixelRGB(100, 100, 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Shift(2)
	}
}

func createPixels(n int) Pixels {
	p := make(Pixels,100)
	for i := 0 ; i < len(p) ; i++ {
		p[i] = NewPixelRGB(uint8(rand.Intn(255)),uint8(rand.Intn(255)),uint8(rand.Intn(255)))
	}
	return p
}
//Pixels
func BenchmarkMost100(b *testing.B) {
	p := createPixels(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Most()
	}
}
func BenchmarkMost10000(b *testing.B) {
	p := createPixels(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Most()
	}
}

func BenchmarkQuantize100(b *testing.B) {
	p := createPixels(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_,err := p.Quantize(2)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkQuantize10000(b *testing.B) {
	p := createPixels(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_,err := p.Quantize(2)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAverage100(b *testing.B) {
	p := createPixels(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_,err := p.Average()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAverage10000(b *testing.B) {
	p := createPixels(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_,err := p.Average()
		if err != nil {
			b.Fatal(err)
		}
	}
}
