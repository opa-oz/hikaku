package hikaku

import (
	"github.com/gkampitakis/go-snaps/snaps"
	"testing"
)

func TestCalcHistogram(t *testing.T) {
	image, err := open("./subjects/golden.jpg")
	second, _ := open("./subjects/copper.jpg")

	histogram := calcHistogram(image, 16)
	secondHist := calcHistogram(second, 16)
	bounds := image.Bounds()

	snaps.MatchSnapshot(t, err)
	snaps.MatchSnapshot(t, histogram)

	resolution := float64(bounds.Max.X * bounds.Max.Y)

	normalizedHistogram := normalizeHistogram(histogram, resolution)
	snaps.MatchSnapshot(t, normalizedHistogram)

	secondNormalizedHistogram := normalizeHistogram(secondHist, resolution)

	red := calcHellinger(normalizedHistogram, secondNormalizedHistogram, 0)   // red
	green := calcHellinger(normalizedHistogram, secondNormalizedHistogram, 1) // green
	blue := calcHellinger(normalizedHistogram, secondNormalizedHistogram, 2)  // blue

	snaps.MatchSnapshot(t, red)
	snaps.MatchSnapshot(t, green)
	snaps.MatchSnapshot(t, blue)

	result := compareHistograms(normalizedHistogram, secondNormalizedHistogram, 16)

	snaps.MatchSnapshot(t, result)
}
