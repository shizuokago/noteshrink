package noteshrink

import (
	"fmt"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"math"
)

func Pack(p *Pixel) int {
	return int(p.R)<<16 | int(p.G)<<8 | int(p.B)
}

func UnPack(v int) (uint8, uint8, uint8) {
	r := uint8((v >> 16) & 0xFF)
	g := uint8((v >> 8) & 0xFF)
	b := uint8(v & 0xFF)
	return r, g, b
}

func convertRGBA(c color.Color) (*color.RGBA, error) {
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

	h := max - min
	if h > 0 {
		if max == r {
			h = (g - b) / h
			if h < 0 {
				h += 6
			}
		} else if max == g {
			h = 2 + (b-r)/h
		} else {
			h = 4 + (r-g)/h
		}
	}
	h /= 6
	s := max - min
	if max > 0 {
		s /= max
	}
	v := max
	return h, s, v
}

//https://www.rapidtables.com/convert/color/hsv-to-rgb.html
func HSV2RGB(h, s, v float64) *color.RGBA {

	r := v
	g := v
	b := v

	if s > 0 {

		h *= 6.
		i := int(h)
		f := h - float64(i)

		switch i {
		default:
		case 0:
			g *= 1 - s*(1-f)
			b *= 1 - s
		case 1:
			r *= 1 - s*f
			b *= 1 - s
		case 2:
			r *= 1 - s
			b *= 1 - s*(1-f)
		case 3:
			r *= 1 - s
			g *= 1 - s*f
		case 4:
			r *= 1 - s*(1-f)
			g *= 1 - s
		case 5:
			g *= 1 - s
			b *= 1 - s*f
		}
	}

	return FloatRGBA(r*255.0, g*255.0, b*255.0)
}

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

func UIntRGBA(r, g, b uint8) *color.RGBA {
	return &color.RGBA{R: r, G: g, B: b, A: 255}
}
