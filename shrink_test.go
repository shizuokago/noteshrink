package noteshrink

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setup()
	ret := m.Run()
	teardown()
	os.Exit(ret)
}

func setup() {
}

func teardown() {
}

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

	samples.output("sample/notesA1_samples.png", 100, 100)

	op.Shift = 4
	bg, err := getBackgroundColor(samples, op)
	if err != nil {
		t.Errorf("BackgroundTest:[%v]", err)
	}

	if bg.R != 224 || bg.G != 224 || bg.B != 224 {
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
	b.ResetTimer()
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

func BenchmarkConvertPixels(b *testing.B) {

	img, err := loadImage("sample/notesA1.jpg")
	if err != nil {
		b.Errorf("loadImage() Error[%v]", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		//データの展開
		_, err := convertPixels(img)
		if err != nil {
			b.Errorf("convertPixel() Error[%v]", err)
		}
	}
}

func BenchmarkCreateSample(b *testing.B) {

	img, err := loadImage("sample/notesA1.jpg")
	if err != nil {
		b.Errorf("loadImage() Error[%v]", err)
		return
	}
	data, err := convertPixels(img)
	if err != nil {
		b.Errorf("convertPixels() Error[%v]", err)
	}

	op := DefaultOption()
	//サンプルの作成
	num := int(float64(len(data)) * op.SamplingRate)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := createSample(data, num)
		if err != nil {
			b.Errorf("createSample() Error[%v]", err)
			return
		}
	}
}

func BenchmarkCreatePalette(b *testing.B) {

	img, err := loadImage("sample/notesA1.jpg")
	if err != nil {
		b.Errorf("loadImage() Error[%v]", err)
		return
	}
	op := DefaultOption()

	//データの展開
	data, err := convertPixels(img)
	if err != nil {
		b.Errorf("convertPixels() Error[%v]", err)
		return
	}

	//サンプルの作成
	num := int(float64(len(data)) * op.SamplingRate)
	samples, err := createSample(data, num)
	if err != nil {
		b.Errorf("createSample() Error[%v]", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		//色の選定
		_, _, err := createPalette(samples, op)
		if err != nil {
			b.Errorf("createPalette() Error[%v]", err)
			return
		}
	}
}

func BenchmarkApply(b *testing.B) {

	img, err := loadImage("sample/notesA1.jpg")
	if err != nil {
		b.Errorf("loadImage() Error[%v]", err)
		return
	}
	op := DefaultOption()

	//データの展開
	data, err := convertPixels(img)
	if err != nil {
		b.Errorf("convertPixels() Error[%v]", err)
	}

	//サンプルの作成
	num := int(float64(len(data)) * op.SamplingRate)
	samples, err := createSample(data, num)
	if err != nil {
		b.Errorf("createSample() Error[%v]", err)
		return
	}

	//色の選定
	bg, palette, err := createPalette(samples, op)
	if err != nil {
		b.Errorf("createPalette() Error[%v]", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		//色の適用
		shrink, err := apply(data, bg, palette, op)
		if err != nil {
			b.Errorf("apply() Error[%v]", err)
			return
		}
		if shrink == nil {
			b.Errorf("apply image is nil")
			return
		}
	}
}

func BenchmarkToImage(b *testing.B) {
	img, err := loadImage("sample/notesA1.jpg")
	if err != nil {
		b.Errorf("loadImage() Error[%v]", err)
		return
	}
	op := DefaultOption()

	//データの展開
	data, err := convertPixels(img)
	if err != nil {
		b.Errorf("convertPixels() Error[%v]", err)
	}

	//サンプルの作成
	num := int(float64(len(data)) * op.SamplingRate)
	samples, err := createSample(data, num)
	if err != nil {
		b.Errorf("createSample() Error[%v]", err)
		return
	}

	//色の選定
	bg, palette, err := createPalette(samples, op)
	if err != nil {
		b.Errorf("createPalett	e() Error[%v]", err)
		return
	}

	//色の適用
	shrink, err := apply(data, bg, palette, op)
	if err != nil {
		b.Errorf("apply() Error[%v]", err)
		return
	}
	if shrink == nil {
		b.Errorf("image is nil")
		return
	}

	rect := img.Bounds()
	cols := rect.Dx()
	rows := rect.Dy()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		shrink.ToImage(cols, rows)
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
