package hikaku

import (
	"github.com/gkampitakis/go-snaps/snaps"
	"testing"
)

func TestMatrix(t *testing.T) {
	golden, _ := open("./subjects/golden.jpg")
	copper, _ := open("./subjects/copper.jpg")

	diffPairs := GetDiffPairs(golden, copper)
	snaps.MatchSnapshot(t, diffPairs)

	diffShapes := GetDiffShapes(golden, copper, ContextParameters{})
	snaps.MatchSnapshot(t, diffShapes)
}
