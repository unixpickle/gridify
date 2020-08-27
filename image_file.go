package main

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/unixpickle/essentials"
)

// ReadImagesFromFile extracts one or more images from a
// file, which may contain a grid of sub-images.
//
// If the filename contains an @, then it will be parsed as
//
//     <path>@<rows>x<cols>
//
// and the images from the grid will be returned in order
// (according to scanline ordering).
func ReadImagesFromFile(name string, deborder bool) (imgs []image.Image, err error) {
	defer essentials.AddCtxTo(fmt.Sprintf("read images %s", name), &err)

	name, rows, cols, err := parseGridFilename(name)
	if err != nil {
		return nil, err
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
	return ExtractFromGrid(img, rows, cols, deborder)
}

func parseGridFilename(name string) (base string, rows, cols int, err error) {
	base = name
	rows = 1
	cols = 1

	idx := strings.LastIndex(name, "@")
	if idx != -1 {
		base = name[:idx]
		dimStr := name[idx+1:]
		parts := strings.Split(dimStr, "x")

		err = fmt.Errorf("invalid grid size after '@': expected <rows>x<cols> but got %s", dimStr)
		if len(parts) != 2 {
			return
		}
		var intErr error
		rows, intErr = strconv.Atoi(parts[0])
		if intErr != nil {
			return
		}
		cols, err = strconv.Atoi(parts[1])
		if intErr != nil {
			return
		}
		err = nil
	}

	return
}

// WriteImageToFile saves an image to a file, using the
// extension as an indicator of the file type.
func WriteImageToFile(outName string, img image.Image) (err error) {
	defer essentials.AddCtxTo(fmt.Sprintf("write image %s", outName), &err)
	w, err := os.Create(outName)
	if err != nil {
		return err
	}
	defer w.Close()

	ext := strings.ToLower(filepath.Ext(outName))
	switch ext {
	case ".png":
		return png.Encode(w, img)
	case ".jpg", ".jpeg":
		return jpeg.Encode(w, img, nil)
	case ".gif":
		return gif.Encode(w, img, nil)
	default:
		return fmt.Errorf("unknown file extension: %s", ext)
	}
}
