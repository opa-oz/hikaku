package hikaku

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
)

var fillColor = color.RGBA{R: 255, G: 0, B: 0, A: 0xff}

type Parameters struct {
	threshold float64
	binsCount int
}

func Compare(golden, copper image.Image, params Parameters) (bool, float64) {
	if params.binsCount == 0 {
		params.binsCount = 16
	}

	if params.threshold == 0.0 {
		params.threshold = 0.2
	}

	isParamsIdentical := CompareByParams(golden, copper)

	if !isParamsIdentical {
		gBounds := golden.Bounds()
		cBounds := copper.Bounds()
		gSquare := float64(gBounds.Max.X * gBounds.Max.Y)
		cSquare := float64(cBounds.Max.X * cBounds.Max.Y)

		biggest := math.Max(gSquare, cSquare)
		diff := math.Abs(gSquare - cSquare)

		return false, diff / biggest
	}

	return CompareByHistograms(golden, copper, params)
}

func CompareByHistograms(golden, copper image.Image, params Parameters) (bool, float64) {
	// https://stackoverflow.com/questions/843972/image-comparison-fast-algorithm
	if params.binsCount == 0 {
		params.binsCount = 16
	}

	if params.threshold == 0.0 {
		params.threshold = 0.2
	}

	bounds := golden.Bounds()

	resolution := float64(bounds.Max.Y * bounds.Max.X)

	normalizedGoldHist := NormalizeHistogram(
		CalcHistogram(golden, params.binsCount),
		resolution,
	)
	normalizedCopperHist := NormalizeHistogram(
		CalcHistogram(copper, params.binsCount),
		resolution,
	)

	result := CompareHistograms(normalizedGoldHist, normalizedCopperHist, params.binsCount)

	fmt.Printf("%f\n", result)

	return result < params.threshold, result
}

func CompareByParams(golden, copper image.Image) bool {
	return golden.Bounds() == copper.Bounds()
}

func getDiffPairs(golden, copper image.Image) []DiffPoint {
	goldenGray := Grayscale(golden)
	copperGray := Grayscale(copper)

	return GetDiffPairs(goldenGray, copperGray)
}

func FindDiffMask(golden, copper image.Image) *image.RGBA {
	bounds := golden.Bounds()
	result := image.NewRGBA(bounds)

	diffPairs := getDiffPairs(golden, copper)

	for i := range diffPairs {
		pair := diffPairs[i]
		result.Set(pair.x, pair.y, fillColor)
	}

	return result
}

func FindDiffShapesMask(golden, copper image.Image) *image.RGBA {
	bounds := golden.Bounds()
	result := image.NewRGBA(bounds)

	shapes := GetDiffShapes(golden, copper)

	for _, rect := range shapes {
		println(rect.minX, rect.minY, rect.maxX, rect.maxY)
		draw.Draw(result, image.Rect(rect.minX, rect.minY, rect.maxX, rect.maxY), &image.Uniform{C: fillColor}, image.Point{}, draw.Src)
	}

	return result
}

func ApplyDiff(golden, diff image.Image) *image.RGBA {
	bounds := golden.Bounds()
	mask := Mask{bounds: bounds, width: bounds.Max.X, height: bounds.Max.Y}

	return mask.Draw(golden, diff, 128)
}
