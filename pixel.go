package noteshrink

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"sort"
)

type Pixel struct {
	R uint8
	G uint8
	B uint8

	H float64
	S float64
	V float64
}

//PackはPixelデータをint化します
func Pack(p *Pixel) int {
	return int(p.R)<<16 | int(p.G)<<8 | int(p.B)
}

//元のRGBデータに直します
func UnPack(v int) (uint8, uint8, uint8) {
	r := uint8((v >> 16) & 0xFF)
	g := uint8((v >> 8) & 0xFF)
	b := uint8(v & 0xFF)
	return r, g, b
}

//ColorからのPixel生成
func NewPixel(c color.Color) *Pixel {
	cc, err := convertColor(c)
	if err != nil {
		return nil
	}
	return NewPixelRGB(cc.R, cc.G, cc.B)
}

//RGBからのPixel生成
func NewPixelRGB(r, g, b uint8) *Pixel {
	p := &Pixel{}
	p.R = r
	p.G = g
	p.B = b
	p.H, p.S, p.V = RGB2HSV(p.R, p.G, p.B)
	return p
}

//HSVからのPixel生成
func NewPixelHSV(h, s, v float64) *Pixel {
	c := HSV2RGB(h, s, v)
	return NewPixelRGB(c.R, c.G, c.B)
}

//HSVの位置を取得
func (p Pixel) DistanceHSV(src *Pixel) (float64, float64, float64) {
	h := math.Abs(src.H - p.H)
	s := math.Abs(src.S - p.S)
	v := math.Abs(src.V - p.V)
	return h, s, v
}

//RGB空間の距離
func (own Pixel) DistanceRGB(src *Pixel) float64 {
	all := 0.0
	r := float64(src.R) - float64(own.R)
	g := float64(src.G) - float64(own.G)
	b := float64(src.B) - float64(own.B)
	all += r * r
	all += g * g
	all += b * b
	return all
}

//Shift変換
func (p Pixel) Shift(shift uint) *Pixel {

	r := uint8((p.R >> shift) << shift)
	g := uint8((p.G >> shift) << shift)
	b := uint8((p.B >> shift) << shift)

	return NewPixelRGB(r, g, b)
}

//色生成
func (p Pixel) Color() *color.RGBA {
	return UIntRGBA(p.R, p.G, p.B)
}

//デバッグ用の文字列作成
func (p Pixel) String() string {
	rtn := fmt.Sprintf("R[%d]G[%d]B[%d] = H[%f]S[%f]V[%f]", p.R, p.G, p.B, p.H, p.S, p.V)
	return rtn
}

type Pixels []*Pixel

//一番多い色を取得
func (p Pixels) Most() *Pixel {

	counter := make(map[int]int)
	for _, pix := range p {
		val := Pack(pix)
		counter[val]++
	}

	max := 0
	value := 0
	for key, elm := range counter {

		if elm > max {
			max = elm
			value = key
		}
	}
	return NewPixelRGB(UnPack(value))
}

//すべてのデータを丸めた色を返す
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

//平均の色を算出
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

//画像の作成
func (p Pixels) ToImage(cols, rows int) (image.Image,error) {

	idx := 0
	img := image.NewRGBA(image.Rect(0, 0, cols, rows))

	for col := 0; col < cols; col++ {
		for row := 0; row < rows; row++ {
			img.Set(col, row, p[idx].Color())
			idx++
		}
	}
	return img,nil
}

//ソート
func (p Pixels) Sort() error {

	sort.Slice(p, func(i, j int) bool {
		pi := Pack(p[i])
		pj := Pack(p[j])
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

//デバッグ用にパレットを作成
func (p Pixels) debug(f string) error {
	leng := len(p)
	if leng > 20 {
		return fmt.Errorf("NotSupported.")
	}
	img,err := p.ToImage(leng, 1)
	if err != nil {
		return err
	}
	return OutputPNG(f, img)
}

//出力
func (p Pixels) output(f string, cols, rows int) error {
	img,err := p.ToImage(cols, rows)
	if err != nil {
		return err
	}
	return OutputPNG(f, img)
}
