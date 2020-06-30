package main

import (
	"fmt"
	"image"

	"github.com/unixpickle/essentials"
)

// PlaceInGrid places all of the tile images into a grid
// in a larger, single image.
//
// If the tiles are not all the same size, smaller tiles
// are padded to be the size of the largest tile.
func PlaceInGrid(tiles []image.Image, cols int) image.Image {
	tileWidth, tileHeight := 0, 0
	for _, img := range tiles {
		bounds := img.Bounds()
		tileWidth = essentials.MaxInt(tileWidth, bounds.Dx())
		tileHeight = essentials.MaxInt(tileHeight, bounds.Dy())
	}

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

// ExtractFromGrid extracts tiles from a grid image.
func ExtractFromGrid(grid image.Image, rows, cols int) ([]image.Image, error) {
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
			results = append(results, &subimage{
				Image:     grid,
				newBounds: image.Rect(x*tileWidth, y*tileHeight, (x+1)*tileWidth, (y+1)*tileHeight),
			})
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
