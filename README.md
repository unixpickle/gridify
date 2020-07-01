# gridify

This is a command-line tool to manipulate image grids. It can read one or more images (possibly extracting input images from their own grids), and then arrange all of the images into a grid in a new image.

```
Usage: gridify [flags] <inputs ...> <output>

Input files may themselves be grids. To read sub-images from a grid,
place an '@' at the end of the name, followed by the grid size:

  <path>@<rows>x<cols> (e.g. /my/image.png@2x4)

Available flags:

  -cols int
        number of columns (0 to layout automatically)
```
