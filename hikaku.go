package hikaku

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

var fillColor = color.RGBA{R: 255, G: 0, B: 0, A: 0xff}

type ComparisonParameters struct {
	Threshold                                float64
	BinsCount                                int
	NormalizedGoldHist, NormalizedCopperHist NormalizedHistogram
}

type ContextParameters struct {
	DiffPairs  []DiffPoint
	DiffShapes map[int]Rect
}

// Compare compares using both - CompareByHistograms and CompareByParams methods
func Compare(golden, copper image.Image, params ComparisonParameters) (bool, float64) {
	if params.BinsCount == 0 {
		params.BinsCount = 16
	}

	if params.Threshold == 0.0 {
		params.Threshold = 0.2
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

// PrepareHistogram prepares histograms for 3 channels (r,g,b) and normalizes them
// It returns normalized histogram for 3 different channels
func PrepareHistogram(image image.Image, params ComparisonParameters) (normalizedHistogram NormalizedHistogram) {
	bounds := image.Bounds()
	resolution := float64(bounds.Max.Y * bounds.Max.X)

	normalizedHistogram = normalizeHistogram(
		calcHistogram(image, params.BinsCount),
		resolution,
	)

	return
}

// CompareByHistograms compares two images by 3 histograms - Red, Green, Blue channels of pixels
// It returns the boolean result and the float64 difference between two images
func CompareByHistograms(golden, copper image.Image, params ComparisonParameters) (bool, float64) {
	// https://stackoverflow.com/questions/843972/image-comparison-fast-algorithm
	if params.BinsCount == 0 {
		params.BinsCount = 16
	}

	if params.Threshold == 0.0 {
		params.Threshold = 0.2
	}

	var normalizedGoldHist, normalizedCopperHist NormalizedHistogram

	if params.NormalizedGoldHist != nil {
		normalizedGoldHist = params.NormalizedGoldHist
	} else {
		normalizedGoldHist = PrepareHistogram(golden, params)
	}

	if params.NormalizedCopperHist != nil {
		normalizedCopperHist = params.NormalizedCopperHist
	} else {
		normalizedCopperHist = PrepareHistogram(copper, params)
	}

	result := compareHistograms(normalizedGoldHist, normalizedCopperHist, params.BinsCount)

	return result < params.Threshold, result
}

// CompareHistogramsOnly is kind of copy of CompareByHistograms, which uses only pre-calculated histograms, instead of images
func CompareHistogramsOnly(params ComparisonParameters) (bool, float64) {
	if params.BinsCount == 0 {
		params.BinsCount = 16
	}

	if params.Threshold == 0.0 {
		params.Threshold = 0.2
	}

	result := compareHistograms(params.NormalizedGoldHist, params.NormalizedCopperHist, params.BinsCount)

	return result < params.Threshold, result
}

// CompareByParams compares width and height of both images
func CompareByParams(golden, copper image.Image) bool {
	return golden.Bounds() == copper.Bounds()
}

// GetDiffPairsGrayscale converts images to gray and gets pairs of pixels, that are different in the both images
// It returns slice of pairs (x,y), representing pixel's coordinates
func GetDiffPairsGrayscale(golden, copper image.Image) []DiffPoint {
	goldenGray := grayscale(golden)
	copperGray := grayscale(copper)

	return GetDiffPairs(goldenGray, copperGray)
}

// FindDiffMask finds pixel-perfect difference between two images
// It returns mask of different pixels
func FindDiffMask(golden, copper image.Image, params ContextParameters) *image.RGBA {
	bounds := golden.Bounds()
	result := image.NewRGBA(bounds)

	var diffPairs []DiffPoint

	if params.DiffPairs != nil {
		diffPairs = params.DiffPairs
	} else {
		diffPairs = GetDiffPairsGrayscale(golden, copper)
	}

	for i := range diffPairs {
		pair := diffPairs[i]
		result.Set(pair.x, pair.y, fillColor)
	}

	return result
}

// FindDiffShapesMask finds shapes (rectangles) inside diff
// It returns mask of highlighted rectangles
func FindDiffShapesMask(golden, copper image.Image, params ContextParameters) *image.RGBA {
	bounds := golden.Bounds()
	result := image.NewRGBA(bounds)

	var diffShapes map[int]Rect

	if params.DiffShapes != nil {
		diffShapes = params.DiffShapes
	} else {
		diffShapes = GetDiffShapes(golden, copper, params)
	}

	for _, rect := range diffShapes {
		if rect.minX == rect.maxX && rect.minY == rect.maxY {
			result.Set(rect.minX, rect.minY, fillColor)
		} else {
			draw.Draw(result, image.Rect(rect.minX, rect.minY, rect.maxX, rect.maxY), &image.Uniform{C: fillColor}, image.Point{}, draw.Src)
		}
	}

	return result
}

// ApplyDiff applies diff as a mask with transparency
// It returns image with highlighted diff
func ApplyDiff(golden, diff image.Image, transparency uint8) *image.RGBA {
	bounds := golden.Bounds()
	mask := Mask{bounds: bounds, width: bounds.Max.X, height: bounds.Max.Y}

	return mask.Draw(golden, diff, transparency)
}
