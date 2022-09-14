package hikaku

import (
	"github.com/gkampitakis/go-snaps/snaps"
	"testing"
)

func TestGrayScaling(t *testing.T) {
	image, err := open("./subjects/golden.jpg")

	image = grayscale(image)

	snaps.MatchSnapshot(t, err)
	snaps.MatchSnapshot(t, image)
}
