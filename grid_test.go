package noteshrink

import (
	"image"
	"os"
)

func loadGrid(f string) (Grid, error) {

	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return ConvertGrid(img)
}
