package main

import (
	"github.com/opa-oz/hikaku"
	"image"
	_ "image/jpeg"
	"image/png"
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

func Save(path string, img image.Image) error {
	pathToFile, err := filepath.Abs(path)
	f, _ := os.Create(pathToFile)
	err = png.Encode(f, img)
	if err != nil {
		return nil
	}

	return nil
}

func main() {
	copperImage, err := Open("./subjects/copper.png")
	goldenImage, err2 := Open("./subjects/golden.png")

	if err != nil {
		println(err.Error())

	}

	if err != nil {
		println(err2.Error())
	}

	byParams := hikaku.CompareByParams(goldenImage, copperImage)

	println("Result by params:", byParams)

	byHistograms, diff := hikaku.CompareByHistograms(goldenImage, copperImage, hikaku.Parameters{})

	println("Result by histograms:", byHistograms, "with diff", diff)

	byBoth, bDiff := hikaku.Compare(goldenImage, copperImage, hikaku.Parameters{})

	println("Result by both:", byBoth, "with diff", bDiff)

	diffMask := hikaku.FindDiffMask(goldenImage, copperImage)
	straightForwardDiff := hikaku.ApplyDiff(goldenImage, diffMask)
	err = Save("./subjects/diff_mask.png", straightForwardDiff)
	if err != nil {
		println(err.Error())
	}

	diffShapes := hikaku.FindDiffShapesMask(goldenImage, copperImage)
	shapesForwardDiff := hikaku.ApplyDiff(goldenImage, diffShapes)
	err = Save("./subjects/diff_shapes.png", shapesForwardDiff)
	if err != nil {
		println(err.Error())
	}
}
