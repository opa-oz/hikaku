package hikaku

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
)

func open(path string) (img image.Image, err error) {
	pathToFile, err := filepath.Abs(path)
	file, err := os.Open(pathToFile)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	img, _, err = image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, err
}

// check checks that indices inside of matrix's bounds
func check(x, y, n, m int) bool {
	return x >= 0 && y >= 0 && x < n-1 && y < m-1
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func max(a, b int) int {
	if a < b {
		return b
	}

	return a
}
