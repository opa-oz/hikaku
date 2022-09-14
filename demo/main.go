package main

import (
	"fmt"
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

func decide(input bool) string {
	if !input {
		return "[Different]"
	}

	return "[Identical]"
}

func makeRound(first, second image.Image, hist1, hist2 hikaku.NormalizedHistogram, title1, title2 string) {
	compParams := hikaku.ComparisonParameters{BinsCount: 16, Threshold: 0.2}

	println("\t", title1, " vs ", title2)
	byParams := hikaku.CompareByParams(first, second)

	println("\t\t1. Compared by params:", decide(byParams))

	println("\t\t2. Compared by histograms:")

	compParams.NormalizedGoldHist = hist1
	compParams.NormalizedCopperHist = hist2
	byHistograms, diff := hikaku.CompareByHistograms(first, second, compParams)

	fmt.Printf("\t\t\t - threshold (0.2): %s (%F)\n", decide(byHistograms), diff)
	compParams.Threshold = 0.1
	byHistograms, diff = hikaku.CompareByHistograms(first, second, compParams)
	fmt.Printf("\t\t\t - threshold (0.1): %s (%F)\n", decide(byHistograms), diff)
	compParams.Threshold = 0.05
	byHistograms, diff = hikaku.CompareByHistograms(first, second, compParams)
	fmt.Printf("\t\t\t - threshold (0.05): %s (%F)\n", decide(byHistograms), diff)

	println("\n\tSave mask diff: subject/" + title1 + "_vs_" + title2 + ".mask.jpg")

	diffPairs := hikaku.GetDiffPairsGrayscale(first, second)
	params := hikaku.ContextParameters{DiffPairs: diffPairs}

	diffMask := hikaku.FindDiffMask(first, second, params)
	straightForwardDiff := hikaku.ApplyDiff(second, diffMask, 128)
	err := Save("./subjects/"+title1+"_vs_"+title2+".mask.jpg", straightForwardDiff)
	if err != nil {
		println(err.Error())
	}

	println("\tSave shapes diff: subject/" + title1 + "_vs_" + title2 + ".shapes.jpg")

	diffShapes := hikaku.FindDiffShapesMask(first, second, params)
	shapesForwardDiff := hikaku.ApplyDiff(second, diffShapes, 128)
	err = Save("./subjects/"+title1+"_vs_"+title2+".shapes.jpg", shapesForwardDiff)
	if err != nil {
		println(err.Error())
	}
}

func main() {
	original, oErr := Open("./subjects/golden.jpg")
	slightlyChanged, sErr := Open("./subjects/copper2.jpg")
	completelyDifferent, cErr := Open("./subjects/copper.jpg")

	if oErr != nil {
		println(oErr.Error())
	}

	if sErr != nil {
		println(sErr.Error())
	}

	if cErr != nil {
		println(sErr.Error())
	}

	compParams := hikaku.ComparisonParameters{BinsCount: 16, Threshold: 0.2}

	originalHist := hikaku.PrepareHistogram(original, compParams)
	slightlyChangedHist := hikaku.PrepareHistogram(slightlyChanged, compParams)
	completelyDifferentHist := hikaku.PrepareHistogram(completelyDifferent, compParams)

	println("# First round")
	makeRound(original, slightlyChanged, originalHist, slightlyChangedHist, "Original", "SlightlyChanged")
	println("\n# Second round")
	makeRound(original, completelyDifferent, originalHist, completelyDifferentHist, "Original", "CompletelyDifferent")

}
