package noteshrink

import (
	"testing"
)

func TestBackground(t *testing.T) {

	prefix := "./cmd/sample/notesA1"
	i := prefix + ".jpg"

	pix, err := loadGrid(i)
	if err != nil {
		t.Errorf("Load pixel[%v]", err)
	}

	op := DefaultOption()
	samples, err := createSample(pix, op)
	if err != nil {
		t.Errorf("CreateSample[%v]", err)
	}

	bg, err := getBackgroundColor(samples, op)
	if err != nil {
		t.Errorf("BackgroundTest:[%v]", err)
	}

	if bg.R != 232 || bg.G != 232 || bg.B != 232 {
		t.Errorf("Background Error:[%v]", bg)
	}

}

func TestQuantaizeSample(t *testing.T) {

	prefix := "./cmd/sample/notesA1"
	i := prefix + ".jpg"
	pix, err := loadGrid(i)
	if err != nil {
		t.Errorf("Load pixel[%v]", err)
	}

	op := DefaultOption()

	samples, err := createSample(pix, op)
	if err != nil {
		t.Errorf("CreateSample[%v]", err)
	}

	q, err := samples.Quantize(2)
	if err != nil {
		t.Errorf("Test Quantize:Pack [%v]", err)
	}
	q.Sort()
	sortFile := prefix + "_quantize.png"

	grid, err := NewGrid(100, 100)
	grid.SetPixels(q)

	err = grid.Output(sortFile)
	if err != nil {
		t.Errorf("Test Sort:Output [%v]", err)
	}
}
