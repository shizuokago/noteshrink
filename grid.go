package noteshrink

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"os"
)

type Grid [][]*Pixel

func NewGrid(rows, cols int) (Grid, error) {
	g := make([][]*Pixel, rows)
	for row := 0; row < rows; row++ {
		g[row] = make([]*Pixel, cols)
	}
	return g, nil
}

func ConvertGrid(img image.Image) (Grid, error) {

	rect := img.Bounds()
	cols := rect.Dx()
	rows := rect.Dy()

	pixels, err := NewGrid(rows, cols)
	if err != nil {
		return nil, err
	}

	for col := 0; col < cols; col++ {
		for row := 0; row < rows; row++ {
			color := img.At(col, row)
			pixels[row][col] = NewPixel(color)
			if err != nil {
				return nil, err
			}
		}
	}
	return pixels, nil
}

func (g Grid) Flat() Pixels {
	cols := g.Cols()
	rows := g.Rows()
	flat := make([]*Pixel, 0, rows*cols)
	for col := 0; col < cols; col++ {
		for row := 0; row < rows; row++ {
			flat = append(flat, g[row][col])
		}
	}
	return flat
}

func (g Grid) SetPixels(p Pixels) (err error) {

	rows := g.Rows()
	cols := g.Cols()

	leng := len(p)

	all := rows * cols
	if len(p) > all {
		err = fmt.Errorf("pixels > grid")
	}

	idx := 0
	for col := 0; col < cols; col++ {
		for row := 0; row < rows; row++ {
			g[row][col] = p[idx]
			idx++
			if leng == idx {
				break
			}
		}
	}
	return err
}

func (p Grid) Rows() int {
	return len(p)
}

func (p Grid) Cols() int {
	return len(p[0])
}

func (p Grid) String() string {
	rows := p.Rows()
	cols := p.Cols()

	var box = bytes.NewBuffer(make([]byte, 0, rows*cols*30))
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			box.WriteString(fmt.Sprintf("%d,%d:[%s]\n", row, col, p[row][col]))
		}
	}
	return box.String()
}

func (g Grid) ToImage() image.Image {
	rows := g.Rows()
	cols := g.Cols()

	img := image.NewRGBA(image.Rect(0, 0, cols, rows))
	for col := 0; col < cols; col++ {
		for row := 0; row < rows; row++ {
			img.Set(col, row, g[row][col].Color())
		}
	}
	return img
}

func (g Grid) Output(f string) error {
	img := g.ToImage()
	out, err := os.Create(f)
	if err != nil {
		return err
	}
	defer out.Close()
	return png.Encode(out, img)
}
