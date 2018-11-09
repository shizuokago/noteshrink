package noteshrink

import (
	"fmt"
	"image/color"
	"math"
	"sort"
)

type Pixel struct {
	value int

	R uint8
	G uint8
	B uint8

	H float64
	S float64
	V float64
}

func NewPixel(c color.Color) *Pixel {

	cc, err := convertRGBA(c)
	if err != nil {
		return nil
	}
	return NewPixelRGB(cc.R, cc.G, cc.B)
}

func NewPixelRGB(r, g, b uint8) *Pixel {
	p := &Pixel{}
	p.R = r
	p.G = g
	p.B = b
	p.H, p.S, p.V = RGB2HSV(p.R, p.G, p.B)
	p.value = Pack(p)
	return p
}

func NewPixelHSV(h, s, v float64) *Pixel {
	c := HSV2RGB(h, s, v)
	return NewPixel(c)
}

func (p Pixel) RGB() (uint8, uint8, uint8) {
	return p.R, p.G, p.B
}

func (p Pixel) DistanceHSV(src *Pixel) (float64, float64, float64) {
	h := math.Abs(src.H - p.H)
	s := math.Abs(src.S - p.S)
	v := math.Abs(src.V - p.V)
	return h, s, v
}

func (own Pixel) DistanceRGB(src *Pixel) float64 {
	all := 0.0
	all += math.Pow(float64(src.R)-float64(own.R), 2)
	all += math.Pow(float64(src.G)-float64(own.G), 2)
	all += math.Pow(float64(src.B)-float64(own.B), 2)
	return all
	//return float64(own.value - src.value)
}

func (p Pixel) Shift(shift uint) *Pixel {

	r := uint8((p.R >> shift) << shift)
	g := uint8((p.G >> shift) << shift)
	b := uint8((p.B >> shift) << shift)

	return NewPixelRGB(r, g, b)
}

func (p Pixel) Color() *color.RGBA {
	return UIntRGBA(p.R, p.G, p.B)
}

func (p Pixel) String() string {
	rtn := fmt.Sprintf("R[%d]G[%d]B[%d] = H[%f]S[%f]V[%f]", p.R, p.G, p.B, p.H, p.S, p.V)
	return rtn
}

type Pixels []*Pixel

func (p Pixels) Most() *Pixel {

	counter := make(map[int]int)
	for _, pix := range p {
		val := pix.value
		counter[val]++
	}

	max := 0
	unpack := 0
	for key, elm := range counter {
		if elm > max {
			max = elm
			unpack = key
		}
	}
	return NewPixelRGB(UnPack(unpack))
}

func (p Pixels) Quantize(s int) (Pixels, error) {

	if s >= 8 {
		return nil, fmt.Errorf("shift not 8 over")
	}

	shift := uint(s)
	quantize := make([]*Pixel, len(p))
	for idx, pix := range p {
		quantize[idx] = pix.Shift(shift)
	}
	return quantize, nil
}

func (p Pixels) Average() (*Pixel, error) {

	if p == nil {
		return nil, fmt.Errorf("Pixels is nil")
	}

	leng := len(p)
	if leng == 0 {
		return nil, fmt.Errorf("Pixels length zero")
	}

	r, g, b := 0.0, 0.0, 0.0

	for _, d := range p {
		r += float64(d.R)
		g += float64(d.G)
		b += float64(d.B)
	}

	ave := 1.0 / float64(leng)
	r = r * ave
	g = g * ave
	b = b * ave

	c := FloatRGBA(r, g, b)
	return NewPixel(c), nil
}

func (p Pixels) Sort() error {

	sort.Slice(p, func(i, j int) bool {
		pi := p[i].value
		pj := p[j].value

		iRGB := int((pi>>16)&0xFF) +
			int((pi>>8)&0xFF) +
			int(pi&0xFF)
		jRGB := int((pj>>16)&0xFF) +
			int((pj>>8)&0xFF) +
			int(pj&0xFF)
		return iRGB < jRGB
	})
	return nil
}

func (p Pixels) Output(f string) error {

	leng := len(p)
	if len(p) > 20 {
		return fmt.Errorf("NotSupported.")
	}

	line := 100
	rows := line
	cols := leng * line

	g, err := NewGrid(rows, cols)
	if err != nil {
		return err
	}

	for col := 0; col < cols; col++ {
		for row := 0; row < rows; row++ {
			idx := col / line
			g[row][col] = p[idx]
		}
	}

	return g.Output(f)
}