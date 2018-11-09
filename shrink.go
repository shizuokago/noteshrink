package noteshrink

import (
	"image"
	"math"
	"math/rand"
	"time"
)

type Option struct {
	SamplingRate    float64
	Brightness      float64
	Saturation      float64
	ForegroundNum   int
	Shift           int
	Iterate         int
	SaturateFlag    bool
	WhiteBackground bool
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

func Shrink(img image.Image, op *Option) (image.Image, error) {

	if op == nil {
		op = DefaultOption()
	}

	grid, err := ConvertGrid(img)
	if err != nil {
		return nil, err
	}

	samples, err := createSample(grid, op)
	if err != nil {
		return nil, err
	}

	bg, palette, err := createPalette(samples, op)
	if err != nil {
		return nil, err
	}

	//Foreground Debug
	//err = palette.Output("./sample/output/notesA1_foreground.png")

	shrink, err := apply(grid, bg, palette, op)
	if err != nil {
		return nil, err
	}

	return shrink.ToImage(), nil
}

func apply(g Grid, bg *Pixel, labels Pixels, op *Option) (Grid, error) {

	rows := g.Rows()
	cols := g.Cols()

	rtn, err := NewGrid(rows, cols)
	if err != nil {
		return nil, err
	}

	flat := g.Flat()
	flag, err := getForegraundMask(flat, bg, op)
	if err != nil {
		return nil, err
	}

	idx := 0
	for col := 0; col < cols; col++ {
		for row := 0; row < rows; row++ {
			newPix := bg
			if flag[idx] {
				wk := closest(flat[idx], labels)
				newPix = labels[wk]
			}
			rtn[row][col] = newPix
			idx++
		}
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

	data := make([]*Pixel, 0, len(p))
	for i, pix := range p {
		if mask[i] {
			data = append(data, pix)
		}
	}

	labels, err := kmeans(data, op)
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

func createSample(g Grid, op *Option) (Pixels, error) {

	orgX := g.Rows()
	orgY := g.Cols()
	num := int(float64(orgX) * float64(orgY) * op.SamplingRate)

	samples := make([]*Pixel, num)
	for idx := 0; idx < num; idx++ {
		x := rand.Intn(orgX)
		y := rand.Intn(orgY)
		samples[idx] = g[x][y]
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
