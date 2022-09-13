package hikaku

import (
	"fmt"
	"image"
)

type Parameters struct {
	threshold float64
	binsCount int
}

func Compare(golden, copper image.Image, params Parameters) bool {
	if params.binsCount == 0 {
		params.binsCount = 16
	}

	if params.threshold == 0.0 {
		params.threshold = 0.2
	}

	return CompareByParams(golden, copper) && CompareByHistograms(golden, copper, params)
}

func CompareByHistograms(golden, copper image.Image, params Parameters) bool {
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

	fmt.Printf("Result: %f \n", 1-result)

	return result >= params.threshold
}

func CompareByParams(golden, copper image.Image) bool {
	return golden.Bounds() == copper.Bounds()
}
