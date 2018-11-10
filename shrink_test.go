package noteshrink

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"testing"
)

func TestBackground(t *testing.T) {

	prefix := "sample/notesA1"
	i := prefix + ".jpg"

	pix, err := loadPixels(i)
	if err != nil {
		t.Errorf("Load pixel[%v]", err)
	}

	//pix.output("sample/notesA1_direct.png",2081,2531)

	op := DefaultOption()

	samples, err := createSample(pix, 10000)
	if err != nil {
		t.Errorf("CreateSample[%v]", err)
	}

	samples.output("sample/notesA1_samples.png",100,100)

	bg, err := getBackgroundColor(samples, op)
	if err != nil {
		t.Errorf("BackgroundTest:[%v]", err)
	}

	if bg.R != 236 || bg.G != 236 || bg.B != 236 {
		t.Errorf("Background Error:[%v]", bg)
	}

}

func TestQuantaizeSample(t *testing.T) {

	prefix := "sample/notesA1"
	i := prefix + ".jpg"
	pix, err := loadPixels(i)
	if err != nil {
		t.Errorf("Load pixel[%v]", err)
		return
	}

	samples, err := createSample(pix, 10000)
	if err != nil {
		t.Errorf("CreateSample[%v]", err)
	}

	q, err := samples.Quantize(2)
	if err != nil {
		t.Errorf("Test Quantize:Pack [%v]", err)
	}
	q.Sort()
	sortFile := prefix + "_quantize.png"

	err = q.output(sortFile, 100, 100)
	if err != nil {
		t.Errorf("Test Sort:Output [%v]", err)
	}
}

func BenchmarkShrink(b *testing.B) {
	img, err := loadImage("sample/notesA1.jpg")
	if err != nil {
		b.Errorf("loadImage() Error[%v]", err)
		return
	}
	op := DefaultOption()
	for i := 0; i < b.N; i++ {
		shrink, err := Shrink(img, op)
		if err != nil {
			b.Errorf("Shrink() Error[%v]", err)
			return
		}
		if shrink == nil {
			b.Errorf("shrink image is nil.")
			return
		}
	}
}

//Test用のツール
func loadImage(f string) (image.Image, error) {

	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

//Test用のツール
func loadPixels(f string) (Pixels, error) {
	img, err := loadImage(f)
	if err != nil {
		return nil, err
	}
	pix, err := convertPixels(img)
	if err != nil {
		return nil, err
	}
	return pix, nil
}
