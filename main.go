package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/unixpickle/essentials"
)

func main() {
	if len(os.Args) < 2 {
		DieUsage()
	}
	subCommand := os.Args[1]
	switch subCommand {
	case "make-grid":
		essentials.Must(MakeGrid(os.Args[2:]))
	default:
		fmt.Fprintln(os.Stderr, "unknown sub-command: "+subCommand)
		fmt.Fprintln(os.Stderr)
		DieUsage()
	}
}

func MakeGrid(args []string) error {
	var cols int
	fs := flag.NewFlagSet("make-grid", flag.ExitOnError)
	fs.IntVar(&cols, "cols", 0, "number of columns (0 means use sqrt)")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: gridify make-grid [flags] <inputs ...> <output>")
		fmt.Fprintln(os.Stderr)
		fs.PrintDefaults()
		fmt.Fprintln(os.Stderr)
		os.Exit(1)
	}
	fs.Parse(args)
	if len(fs.Args()) < 2 {
		fs.Usage()
	}

	var tiles []image.Image
	for _, imageName := range fs.Args()[:len(fs.Args())-1] {
		images, err := ImagesFromFilename(imageName)
		if err != nil {
			return fmt.Errorf("read %s: %s", imageName, err.Error())
		}
		tiles = append(tiles, images...)
	}
	if cols == 0 {
		cols = int(math.Ceil(math.Sqrt(float64(len(tiles)))))
		if cols == 0 {
			cols = 1
		}
	}

	output := PlaceInGrid(tiles, cols)
	return WriteImageToFile(fs.Args()[len(fs.Args())-1], output)
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

func DieUsage() {
	fmt.Fprintln(os.Stderr, "Usage: gridify sub_command args")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Available sub-commands:")
	fmt.Fprintln(os.Stderr, "  make-grid [-cols N] <inputs ...> <output>")
	fmt.Fprintln(os.Stderr, "")
	os.Exit(1)
}
