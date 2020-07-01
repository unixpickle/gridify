package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/unixpickle/essentials"
)

func main() {
	var cols int
	flag.IntVar(&cols, "cols", 0, "number of columns (0 to layout automatically)")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: gridify [flags] <inputs ...> <output>")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Input files may themselves be grids. To read sub-images from a grid,")
		fmt.Fprintln(os.Stderr, "place an '@' at the end of the name, followed by the grid size:")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "  <path>@<rows>x<cols> (e.g. /my/image.png@2x4)")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Available flags:")
		fmt.Fprintln(os.Stderr)
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr)
		os.Exit(1)
	}
	flag.Parse()
	if len(flag.Args()) < 2 {
		flag.Usage()
	}

	var tiles []image.Image
	for _, imageName := range flag.Args()[:len(flag.Args())-1] {
		images, err := ImagesFromFilename(imageName)
		if err != nil {
			essentials.Die(fmt.Sprintf("read %s: %s", imageName, err.Error()))
		}
		tiles = append(tiles, images...)
	}
	if cols == 0 {
		cols = AutoGridColumns(tiles)
	}

	output := PlaceInGrid(tiles, cols)
	essentials.Must(WriteImageToFile(flag.Args()[len(flag.Args())-1], output))
}

func ImagesFromFilename(name string) ([]image.Image, error) {
	idx := strings.LastIndex(name, "@")
	rows := 1
	cols := 1
	if idx != -1 {
		dimensions := name[idx+1:]
		name = name[:idx]
		parts := strings.Split(dimensions, "x")
		partErr := fmt.Errorf("invalid grid size string after '@': %s (expected NxM)", dimensions)
		if len(parts) != 2 {
			return nil, partErr
		}
		var err error
		rows, err = strconv.Atoi(parts[0])
		if err != nil {
			return nil, partErr
		}
		cols, err = strconv.Atoi(parts[1])
		if err != nil {
			return nil, partErr
		}
	}
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return ExtractFromGrid(img, rows, cols)
}

func WriteImageToFile(outName string, img image.Image) error {
	w, err := os.Create(outName)
	if err != nil {
		return err
	}
	defer w.Close()
	if filepath.Ext(outName) == ".png" {
		return png.Encode(w, img)
	} else {
		return jpeg.Encode(w, img, nil)
	}
}
