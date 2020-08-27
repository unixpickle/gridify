package main

import (
	"fmt"
	"image"

	"github.com/unixpickle/essentials"
)

// AutoGridColumns automatically determines a good number
// of columns for laying out tiles in a grid.
func AutoGridColumns(tiles []image.Image) int {
	tileWidth, tileHeight := tileBounds(tiles)

	var bestCols int
	var bestCircumference int
	for cols := 1; cols <= len(tiles); cols++ {
		rows := len(tiles) / cols
		if len(tiles)%cols != 0 {
			rows++
		}
		circum := tileWidth*cols + tileHeight*rows
		// Use <= here to get as wide an image as possible.
		if cols == 1 || circum <= bestCircumference {
			bestCols = cols
			bestCircumference = circum
		}
	}

	return bestCols
}

// PlaceInGrid places all of the tile images into a grid
// in a larger, single image.
//
// If the tiles are not all the same size, smaller tiles
// are padded to be the size of the largest tile.
//
// The border argument specifies the number of extra
// pixels to put around the largest tile.
func PlaceInGrid(tiles []image.Image, cols, border int) image.Image {
	tileWidth, tileHeight := tileBounds(tiles)
	tileWidth += border
	tileHeight += border

	rows := len(tiles) / cols
	if len(tiles)%cols != 0 {
		rows++
	}

	result := image.NewRGBA(image.Rect(0, 0, tileWidth*cols, tileHeight*rows))
	for i, img := range tiles {
		row := i / cols
		col := i % cols
		outX := col * tileWidth
		outY := row * tileHeight
		bounds := img.Bounds()
		outX += (tileWidth - bounds.Dx()) / 2
		outY += (tileHeight - bounds.Dy()) / 2

		for y := 0; y < bounds.Dy(); y++ {
			for x := 0; x < bounds.Dx(); x++ {
				src := img.At(x+bounds.Min.X, y+bounds.Min.Y)
				result.Set(outX+x, outY+y, src)
			}
		}
	}

	return result
}

func tileBounds(tiles []image.Image) (width, height int) {
	for _, img := range tiles {
		bounds := img.Bounds()
		width = essentials.MaxInt(width, bounds.Dx())
		height = essentials.MaxInt(height, bounds.Dy())
	}
	return
}

// ExtractFromGrid extracts tiles from a grid image.
func ExtractFromGrid(grid image.Image, rows, cols int, deborder bool) ([]image.Image, error) {
	if grid.Bounds().Dx()%cols != 0 {
		return nil, fmt.Errorf("number of columns (%d) does not divide width (%d)",
			cols, grid.Bounds().Dx())
	} else if grid.Bounds().Dy()%rows != 0 {
		return nil, fmt.Errorf("number of rows (%d) does not divide height (%d)",
			cols, grid.Bounds().Dy())
	}

	tileWidth := grid.Bounds().Dx() / cols
	tileHeight := grid.Bounds().Dy() / rows

	var results []image.Image
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			si := &subimage{
				Image:     grid,
				newBounds: image.Rect(x*tileWidth, y*tileHeight, (x+1)*tileWidth, (y+1)*tileHeight),
			}
			if deborder {
				si.RemoveBorder()
			}
			results = append(results, si)
		}
	}

	return results, nil
}

type subimage struct {
	image.Image
	newBounds image.Rectangle
}

func (s *subimage) Bounds() image.Rectangle {
	return s.newBounds
}

func (s *subimage) RemoveBorder() {
	minX := s.newBounds.Max.X
	minY := s.newBounds.Max.Y
	maxX := s.newBounds.Min.X
	maxY := s.newBounds.Min.Y
	for y := s.newBounds.Min.Y; y < s.newBounds.Max.Y; y++ {
		for x := s.newBounds.Min.X; x < s.newBounds.Max.X; x++ {
			_, _, _, a := s.Image.At(x, y).RGBA()
			if a > 0 {
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	// Handle fully-empty images.
	maxX = essentials.MaxInt(maxX, minX)
	maxY = essentials.MaxInt(maxY, minY)

	s.newBounds.Min.X = minX
	s.newBounds.Min.Y = minY
	s.newBounds.Max.X = maxX
	s.newBounds.Max.Y = maxY
}
