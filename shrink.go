/*
 Shrink() を呼び出すと image.Image をnoteshrinkして image.Imageに変換してくれます。


*/
package noteshrink

import (
	"fmt"
	"image"
	"math"
	"math/rand"
	"time"
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

	//データの展開
	data, err := convertPixels(img)
	if err != nil {
		return nil, err
	}

	//サンプルの作成
	num := int(float64(len(data)) * op.SamplingRate)
	samples, err := createSample(data, num)
	if err != nil {
		return nil, err
	}

	//色の選定
	bg, palette, err := createPalette(samples, op)
	if err != nil {
		return nil, err
	}

	//色の適用
	shrink, err := apply(data, bg, palette, op)
	if err != nil {
		return nil, err
	}
	if shrink == nil {
		return nil, fmt.Errorf("image is nil")
	}

	//GIF用にパレットを作成
	setGIFPalette(bg, palette)

	rect := img.Bounds()
	cols := rect.Dx()
	rows := rect.Dy()

	return shrink.ToImage(cols, rows), nil
}

//色を適用
func apply(data Pixels, bg *Pixel, labels Pixels, op *Option) (Pixels, error) {

	//使用箇所を取得
	flag, err := getForegraundMask(data, bg, op)
	if err != nil {
		return nil, err
	}

	rtn := make([]*Pixel, len(data))
	for idx := 0; idx < len(data); idx++ {
		newPix := bg
		if flag[idx] {
			//近いラベルを取得
			wk := closest(data[idx], labels)
			newPix = labels[wk]
		}
		rtn[idx] = newPix
	}
	return rtn, nil
}

//使用する色を検索
func createPalette(p Pixels, op *Option) (*Pixel, Pixels, error) {

	//背景色を取得
	bg, err := getBackgroundColor(p, op)
	if err != nil {
		return nil, nil, err
	}

	//使用箇所を特定
	mask, err := getForegraundMask(p, bg, op)
	if err != nil {
		return bg, nil, err
	}

	//適用だけを残す
	target := make([]*Pixel, 0, len(p))
	for i, pix := range p {
		if mask[i] {
			target = append(target, pix)
		}
	}

	//色を決定
	labels, err := kmeans(target, op)
	if err != nil {
		return bg, nil, err
	}

	return bg, labels, nil
}

//背景色を取得
func getBackgroundColor(p Pixels, op *Option) (*Pixel, error) {

	//色を落とす
	q, err := p.Quantize(op.Shift)
	if err != nil {
		return nil, err
	}
	//一番多い色を取得
	col := q.Most()

	return col, nil
}

//サンプルを抽出
func createSample(p Pixels, num int) (Pixels, error) {

	samples := make([]*Pixel, num)
	leng := len(p)
	for idx := 0; idx < num; idx++ {
		samples[idx] = p[rand.Intn(leng)]
	}
	return samples, nil
}

//HSV空間からの距離により、使用箇所を特定
func getForegraundMask(p Pixels, bg *Pixel, op *Option) ([]bool, error) {

	rtn := make([]bool, len(p))
	for idx, pix := range p {
		_, ds, dv := pix.DistanceHSV(bg)
		rtn[idx] = dv >= op.Brightness || ds >= op.Saturation
	}
	return rtn, nil
}

//kmeansで色を特定
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
				//TODO error
				//return nil,err
			}
		}

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

//近い位置を取得
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


//一般化してみたやつ（未使用）
type Value interface {
	Distance(Value) float64
	Average([]Value) (Value, error)
}

func kmeansValue(data []Value, labels []Value, itr int) []Value {

	index := make([]int, len(data))
	for idx, datum := range data {
		index[idx] = closestIndex(datum, labels)
	}

	rtn := make([]Value, len(labels))
	for idx, label := range labels {
		rtn[idx] = label
	}

	for idx := 0; idx < itr; idx++ {

		groups := make([][]Value, len(rtn))
		for i := range rtn {
			groups[i] = make([]Value, 0, len(data))
		}

		for i, elm := range data {
			idx := index[i]
			groups[idx] = append(groups[idx], elm)
		}

		for i, label := range rtn {
			valSlice := groups[i]
			ave, err := label.Average(valSlice)
			if ave != nil && err == nil {
				rtn[i] = ave
			} else if err != nil {
			}
		}

		changes := 0
		for i, pix := range data {
			if newIdx := closestIndex(pix, rtn); newIdx != index[i] {
				changes++
				index[i] = newIdx
			}
		}

		if changes == 0 {
			break
		}
	}
	return rtn
}

func closestIndex(val Value, labels []Value) int {
	rtn := -1
	min := math.MaxFloat64
	for idx, elm := range labels {
		vd := val.Distance(elm)
		if vd < min {
			min = vd
			rtn = idx
		}
	}
	return rtn
}
