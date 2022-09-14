package hikaku

import (
	"image"
	"math"
)

type Histogram = [][3]int
type NormalizedHistogram = [][3]float64

func calcHistogram(image image.Image, bins int, maxOpt ...int) Histogram {
	// http://www.sci.utah.edu/~acoste/uou/Image/project1/Arthur_COSTE_Project_1_report.html
	max := 65535

	if len(maxOpt) > 0 {
		max = maxOpt[0]
	}

	histogram := make(Histogram, bins)

	binSize := max / bins

	bounds := image.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := image.At(x, y).RGBA()

			for binIndex := 0; binIndex < bins; binIndex++ {
				low := uint32(binIndex * binSize)
				high := uint32((binIndex - 1) * binSize)

				if r <= low && r > high {
					histogram[binIndex][0]++
				}
				if g <= low && g > high {
					histogram[binIndex][1]++
				}
				if b <= low && b > high {
					histogram[binIndex][2]++
				}
			}
		}
	}

	return histogram
}

func normalizeHistogram(histogram Histogram, resolution float64) NormalizedHistogram {
	newHist := make(NormalizedHistogram, len(histogram))

	for i := range histogram {
		newHist[i][0] = float64(histogram[i][0]) / resolution
		newHist[i][1] = float64(histogram[i][1]) / resolution
		newHist[i][2] = float64(histogram[i][2]) / resolution
	}

	return newHist
}

func calcHellinger(golden NormalizedHistogram, copper NormalizedHistogram, index int) float64 {
	//https://en.wikipedia.org/wiki/Hellinger_distance
	result := 0.0

	for i := range golden {
		result += math.Pow(math.Sqrt(golden[i][index])-math.Sqrt(copper[i][index]), 2)
	}

	result = 1 / math.Sqrt(2) * math.Sqrt(result)

	return result
}

func compareHistograms(golden NormalizedHistogram, copper NormalizedHistogram, bins int) float64 {
	goldenRed := make([]float64, bins)
	copperRed := make([]float64, bins)

	for i := range golden {
		goldenRed[i] = golden[i][0]
		copperRed[i] = copper[i][0]
	}

	red := calcHellinger(golden, copper, 0)   // red
	green := calcHellinger(golden, copper, 1) // green
	blue := calcHellinger(golden, copper, 2)  // blue

	return red + green + blue
}
