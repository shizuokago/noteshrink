package noteshrink

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"math"
	"os"
)

//PNG の圧縮出力
func OutputPNG(f string, img image.Image) error {
	//出力ファイルの作成
	out, err := os.Create(f)
	if err != nil {
		return err
	}
	defer out.Close()

	var enc png.Encoder
	enc.CompressionLevel = png.BestCompression
	return enc.Encode(out, img)
}

var gifPalette color.Palette = nil

//減色GIFのパレット作成
func setGIFPalette(bg *Pixel, fore Pixels) {
	gifPalette = make(color.Palette, len(fore)+1)
	gifPalette[0] = bg.Color()
	for i, pix := range fore {
		gifPalette[i+1] = pix.Color()
	}
}

//減色したGIFパレットでの出力
func OutputGIF(f string, img image.Image) error {

	if gifPalette == nil {
		return fmt.Errorf("palette is nil")
	}

	//出力ファイルの作成
	out, err := os.Create(f)
	if err != nil {
		return err
	}
	defer out.Close()

	op := &gif.Options{
		NumColors: len(gifPalette),
		Quantizer: NewQuantizer(gifPalette),
	}
	return gif.Encode(out, img, op)
}

//減色GIFのQuantazer
type gifQuantizer struct {
	palette color.Palette
}

//Quantizerの生成
func NewQuantizer(p color.Palette) *gifQuantizer {
	q := gifQuantizer{}
	q.palette = p
	return &q
}

//Quantizer実装
func (q gifQuantizer) Quantize(p color.Palette, img image.Image) color.Palette {
	return q.palette
}

//ConvertGridにより、image.ImageをGridに展開します
func convertPixels(img image.Image) (Pixels, error) {

	rect := img.Bounds()
	cols := rect.Dx()
	rows := rect.Dy()

	rtn := make(Pixels, cols*rows)
	idx := 0

	for col := 0; col < cols; col++ {
		for row := 0; row < rows; row++ {
			color := img.At(col, row)
			rtn[idx] = NewPixel(color)
			idx++
		}
	}

	return rtn, nil
}

//colorのキャスト
func convertColor(c color.Color) (*color.RGBA, error) {
	switch c.(type) {
	case color.YCbCr:
		o := c.(color.YCbCr)
		r, g, b := color.YCbCrToRGB(o.Y, o.Cb, o.Cr)
		return UIntRGBA(r, g, b), nil
	case color.RGBA:
		newColor := c.(color.RGBA)
		return &newColor, nil
	case *color.RGBA:
		newColor := c.(*color.RGBA)
		return newColor, nil
	default:
	}
	return nil, fmt.Errorf("not support color[%v]", c)
}

//https://www.rapidtables.com/convert/color/rgb-to-hsv.html
func RGB2HSV(or, og, ob uint8) (float64, float64, float64) {

	r := float64(or) / float64(255)
	g := float64(og) / float64(255)
	b := float64(ob) / float64(255)

	max := math.Max(math.Max(r, g), b)
	min := math.Min(math.Min(r, g), b)

	d := max - min
	h := 0.0
	switch {
	case d == 0:
		h = 0
	case max == r:
		h = math.Mod((g-b)/d, 6)
	case max == g:
		h = (b-r)/d + 2
	case max == b:
		h = (r-g)/d + 4
	}
	h = h / 6
	if h < 0 {
		h += 1.0
	}

	s := 0.0
	if max != 0 {
		s = d / max
	}

	v := max

	return h, s, v
}

//https://www.rapidtables.com/convert/color/hsv-to-rgb.html
func HSV2RGB(h, s, v float64) *color.RGBA {

	hd := h * 360.0
	if hd >= 360 {
		hd = 359
	}

	hh := hd / 60

	c := v * s
	x := c * (1.0 - math.Abs(math.Mod(hh, 2)-1.0))

	r := 0.0
	g := 0.0
	b := 0.0

	switch {
	case hh < 1:
		r, g, b = c, x, 0
	case hh < 2:
		r, g, b = x, c, 0
	case hh < 3:
		r, g, b = 0, c, x
	case hh < 4:
		r, g, b = 0, x, c
	case hh < 5:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	m := v - c
	r += m
	g += m
	b += m

	return FloatRGBA(r*255.0, g*255.0, b*255.0)
}

//FloatのRGB値からRGBAの作成
func FloatRGBA(r, g, b float64) *color.RGBA {

	ur := uint8(r)
	ug := uint8(g)
	ub := uint8(b)
	if r < 255.0 {
		ur = uint8(math.Floor(r + 0.5))
	}
	if g < 255.0 {
		ug = uint8(math.Floor(g + 0.5))
	}
	if b < 255.0 {
		ub = uint8(math.Floor(b + 0.5))
	}
	return UIntRGBA(ur, ug, ub)
}

//RGBAの作成
func UIntRGBA(r, g, b uint8) *color.RGBA {
	return &color.RGBA{R: r, G: g, B: b, A: 255}
}
