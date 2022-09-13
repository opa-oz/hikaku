package main

import (
	"github.com/opa-oz/hikaku"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
)

func Open(path string) (img image.Image, err error) {
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

func main() {
	copperImage, err := Open("./subjects/copper.jpg")
	goldenImage, err2 := Open("./subjects/golden.jpg")

	if err != nil {
		println(err.Error())

	}

	if err != nil {
		println(err2.Error())
	}

	byParams := hikaku.CompareByParams(goldenImage, copperImage)

	println("Result by params:", byParams)

	byHistograms := hikaku.CompareByHistograms(goldenImage, copperImage, hikaku.Parameters{})

	println("Result by histograms:", byHistograms)

	byBoth := hikaku.Compare(goldenImage, copperImage, hikaku.Parameters{})

	println("Result by both:", byBoth)
}
