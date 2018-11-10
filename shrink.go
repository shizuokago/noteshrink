/*
 Shrink を呼び出すと image.Image をnoteshrinkして image.Imageに変換してくれます。


*/
package noteshrink

import (
	"fmt"
	"image"
	"math"
	"math/rand"
	"time"
	_ "image/jpeg"
	_ "image/png"
)

//Option はロジックに対し
type Option struct {
	SamplingRate  float64
	Brightness    float64
	Saturation    float64
	ForegroundNum int
	Shift         int
	Iterate       int
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func DefaultOption() *Option {
	return &Option{
		SamplingRate:  0.002,
		Brightness:    0.30,
		Saturation:    0.20,
		Shift:         2,
		ForegroundNum: 6,
		Iterate:       40,
	}
}

//圧縮
func Shrink(img image.Image, op *Option) (image.Image, error) {

	if op == nil {
		op = DefaultOption()
	}

	data, err := convertPixels(img)
	if err != nil {
		return nil, err
	}

	num := int(float64(len(data)) * op.SamplingRate)
	samples, err := createSample(data, num)
	if err != nil {
		return nil, err
	}

	bg, palette, err := createPalette(samples, op)
	if err != nil {
		return nil, err
	}

	shrink, err := apply(data, bg, palette, op)
	if err != nil {
		return nil, err
	}
	if shrink == nil {
		return nil, fmt.Errorf("image is nil")
	}

	rect := img.Bounds()
	cols := rect.Dx()
	rows := rect.Dy()

	return shrink.ToImage(cols, rows), nil
}

func apply(data Pixels, bg *Pixel, labels Pixels, op *Option) (Pixels, error) {

	flag, err := getForegraundMask(data, bg, op)
	if err != nil {
		return nil, err
	}

	rtn := make([]*Pixel, len(data))
	for idx := 0; idx < len(data); idx++ {
		newPix := bg
		if flag[idx] {
			wk := closest(data[idx], labels)
			newPix = labels[wk]
		}
		rtn[idx] = newPix
	}
	return rtn, nil
}

func createPalette(p Pixels, op *Option) (*Pixel, Pixels, error) {

	bg, err := getBackgroundColor(p, op)
	if err != nil {
		return nil, nil, err
	}

	mask, err := getForegraundMask(p, bg, op)
	if err != nil {
		return bg, nil, err
	}

	target := make([]*Pixel, 0, len(p))
	for i, pix := range p {
		if mask[i] {
			target = append(target, pix)
		}
	}

	labels, err := kmeans(target, op)
	if err != nil {
		return bg, nil, err
	}

	return bg, labels, nil
}

func getBackgroundColor(p Pixels, op *Option) (*Pixel, error) {

	q, err := p.Quantize(op.Shift)
	if err != nil {
		return nil, err
	}
	col := q.Most()

	return col, nil
}

func createSample(p Pixels, num int) (Pixels, error) {

	samples := make([]*Pixel, num)
	leng := len(p)
	for idx := 0; idx < num; idx++ {
		samples[idx] = NewPixel(p[rand.Intn(leng)].Color())
	}

	return samples, nil
}

func getForegraundMask(p Pixels, bg *Pixel, op *Option) ([]bool, error) {

	rtn := make([]bool, len(p))
	for idx, pix := range p {
		_, ds, dv := pix.DistanceHSV(bg)
		rtn[idx] = dv >= op.Brightness || ds >= op.Saturation
	}
	return rtn, nil
}

func kmeans(p Pixels, op *Option) ([]*Pixel, error) {

	k := op.ForegroundNum - 1
	itr := op.Iterate

	labels := make([]*Pixel, k)
	for i := 0; i < k; i++ {
		h := float64(i) / float64(k-1)
		pixel := NewPixelHSV(h, 1, 1)
		labels[i] = pixel
	}

	index := make([]int, len(p))
	for idx, pix := range p {
		index[idx] = closest(pix, labels)
	}

	for idx := 0; idx < itr; idx++ {

		//TODO Routine

		groups := make([]Pixels, len(labels))
		for i := range labels {
			groups[i] = make([]*Pixel, 0, len(labels))
		}

		for i, pix := range p {
			label := index[i]
			groups[label] = append(groups[label], pix)
		}

		for i := range labels {
			if newLabel, err := groups[i].Average(); newLabel != nil && err == nil {
				labels[i] = newLabel
			} else if err != nil {
				//TODO エラー
			}
		}

		//TODO routine end

		changes := 0
		for i, pix := range p {
			if newIdx := closest(pix, labels); newIdx != index[i] {
				changes++
				index[i] = newIdx
			}
		}

		if changes == 0 {
			break
		}
	}

	return labels, nil
}

func closest(p *Pixel, labels []*Pixel) int {

	idx := -1
	d := math.MaxFloat64
	for i := 0; i < len(labels); i++ {
		val := p.DistanceRGB(labels[i])
		if val < d {
			d = val
			idx = i
		}
	}
	return idx
}
