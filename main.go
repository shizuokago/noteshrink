package noteshrink

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
	"strings"
)

var (
	samplingRateOpt = flag.Float64("r", 0.002, "背景色、前景色を選定する際のサンプル数の割合。")
	shiftOpt        = flag.Int("shift", 2, "画素圧縮時のシフト数")

	brigthnessOpt = flag.Float64("b", 0.3, "前景色選定時のVの距離")
	saturationOpt = flag.Float64("s", 0.2, "前景色選定時のSの距離")

	foregroundNumOpt = flag.Int("f", 6, "前景色に選ばれる数を指定")
	iterateOpt       = flag.Int("i", 40, "kmeans が探索するループ数")
)

func Usage() {
	fmt.Println("引数は変換するファイルを複数指定できます")
	flag.Usage()
}

//https://mzucker.github.io/2016/09/20/noteshrink.html
func main() {

	//flagを処理
	flag.Parse()
	opt := Option{
		SamplingRate:  *samplingRateOpt,
		Shift:         *shiftOpt,
		Brightness:    *brigthnessOpt,
		Saturation:    *saturationOpt,
		ForegroundNum: *foregroundNumOpt,
		Iterate:       *iterateOpt,
	}

	files := flag.Args()

	if files == nil || len(files) == 0 {
		Usage()
		return
	}

	for _, f := range files {
		err := run(f, &opt)
		if err != nil {
			fmt.Printf("[%v]\n", err)
			os.Exit(1)
		}
	}

	os.Exit(0)
	return
}

func run(f string, opt *Option) error {

	output := ""
	suffix := "_min"
	idx := strings.LastIndex(f, ".")
	if idx == -1 {
		output = f + suffix
	} else {
		output = f[:idx] + suffix + f[idx:]
	}

	//画像の読み込み
	img, err := loadImage(f)
	if err != nil {
		return err
	}

	//圧縮
	shrink, err := Shrink(img, opt)
	if err != nil {
		return err
	}
	if img == nil {
		return fmt.Errorf("shrink image is null.")
	}

	//出力ファイルの作成
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()

	var enc png.Encoder
	enc.CompressionLevel = png.BestCompression
	return enc.Encode(out, shrink)
}

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
