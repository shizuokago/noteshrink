package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"strings"
	"runtime/pprof"
	"sync"

	"github.com/shizuokago/noteshrink"
)

var (
	samplingRateOpt = flag.Float64("r", 0.002, "背景色、前景色を選定する際のサンプル数の割合。")
	shiftOpt        = flag.Int("shift", 2, "画素圧縮時のシフト数")

	brightnessOpt = flag.Float64("b", 0.35, "前景色選定時のVの距離")
	saturationOpt = flag.Float64("s", 0.25, "前景色選定時のSの距離")

	foregroundNumOpt = flag.Int("f", 6, "前景色に選ばれる数を指定")
	iterateOpt       = flag.Int("i", 40, "kmeans のループ数")

	profileVal = flag.String("p", "", "プロファイル名（指定しない場合プロファイルを行わない）")
	suffixVal  = flag.String("suffix", "_min", "変換ファイル名のサフィックス")
	gifVal     = flag.Bool("g", false, "GIF化したもの")
)

func Usage() {
	fmt.Println("引数は変換するファイルを複数指定できます")
	flag.Usage()
}

//https://mzucker.github.io/2016/09/20/noteshrink.html
func main() {

	//flagを処理
	flag.Parse()
	//プロファイリングを行う
	if *profileVal != "" {
		defer startProfile(*profileVal).stop()
	}

	//オプションをflagから設定
	opt := noteshrink.Option{
		SamplingRate:  *samplingRateOpt,
		Shift:         *shiftOpt,
		Brightness:    *brightnessOpt,
		Saturation:    *saturationOpt,
		ForegroundNum: *foregroundNumOpt,
		Iterate:       *iterateOpt,
	}

	//ファイル名を処理する
	files := flag.Args()
	if files == nil || len(files) == 0 {
		Usage()
		return
	}

	//各処理を非同期で行う
	wg := sync.WaitGroup{}
	for _, f := range files {
		wg.Add(1)
		go func(file string) {
			err := run(file, &opt)
			if err != nil {
				fmt.Printf("[%v]\n", err)
			}
			wg.Done()
		}(f)
	}
	wg.Wait()

	return
}

//ファイル変換の実行
func run(f string, opt *noteshrink.Option) error {

	log.Printf("Shrink    : [%s]\n", f)

	//画像の読み込み
	img, err := loadImage(f)
	if err != nil {
		return err
	}

	//圧縮
	shrink, err := noteshrink.Shrink(img, opt)
	if err != nil {
		return err
	}
	if shrink == nil {
		return fmt.Errorf("shrink image is null.")
	}

	output := ""
	ext := ".png"
	if *gifVal {
		ext = ".gif"
	}
	idx := strings.LastIndex(f, ".")
	//出力ファイル名
	if idx == -1 {
		output = f + *suffixVal + ext
	} else {
		output = f[:idx] + *suffixVal + ext
	}

	//出力の切り替え
	if *gifVal {
		err = noteshrink.OutputGIF(output, shrink)
	} else {
		err = noteshrink.OutputPNG(output, shrink)
	}

	if err == nil {
		log.Printf("Generated : [%s]\n", output)
	}

	return err
}

//画像の読み込み
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

type profile struct {
	file *os.File
	err  error
}

//負荷計測用のプロファイル作成
func startProfile(f string) *profile {
	log.Println("Profile Start:" + f)
	rtn := profile{}
	file, err := os.Create(f)
	if err != nil {
		rtn.err = err
	} else {
		err = pprof.StartCPUProfile(file)
		if err == nil {
			rtn.file = file
		} else {
			rtn.err = err
			defer file.Close()
		}
	}
	return &rtn
}

//プロファイルの終了
func (p profile) stop() {
	log.Println("Profile Stop")
	if p.err == nil {
		pprof.StopCPUProfile()
		p.file.Close()
	} else {
		log.Println(p.err)
	}
}
