package main

import (
	"flag"
	"fmt"
	"image"
	"os"

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
		images, err := ReadImagesFromFile(imageName)
		essentials.Must(err)
		tiles = append(tiles, images...)
	}
	if cols == 0 {
		cols = AutoGridColumns(tiles)
	}

	output := PlaceInGrid(tiles, cols)
	essentials.Must(WriteImageToFile(flag.Args()[len(flag.Args())-1], output))
}
