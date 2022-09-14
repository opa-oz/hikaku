# Hikaku

<p align="center">Yet another tool for image comparison.</p> 

<p align="center">
<img src="https://github.com/opa-oz/hikaku/raw/main/docs/golden.png" width="150" alt="Original"/>
<img src="https://github.com/opa-oz/hikaku/raw/main/docs/copper.png" width="150" alt="Changed"/>
</p>
<p align="center">
<img src="https://github.com/opa-oz/hikaku/raw/main/docs/mask.png" width="150" alt="Changed"/>
<br>
<img src="https://github.com/opa-oz/hikaku/raw/main/docs/shape.png" width="150" alt="Changed"/><sub>[*]</sub>
</p>

## 比較 [ひかく]

_[en: hikaku](https://www.nihongomaster.com/japanese/dictionary/word/46485/hikaku-比較-ひかく)_

Verb
<br>比較

to compare; to contras

## Motivation

**Hikaku** was created for my projects, focused on generating images. I **really** want to compare generated images or
screenshots **quickly** and **reliably**.

I'll mention examples below:

| Repository | Purpose | Example |
|------------|---------|---------|
| tbd        |         |         |
| tbd        |         |         |
| tbd        |         |         |

## Usage

### Compare two images

```go
package main

import (
	"github.com/opa-oz/hikaku"
	"image"
)

func compareMyImages(first, second image.Image) {
	// compare width and height of images
	isBoundsEqual := hikaku.CompareByParams(first, second)

	// build histogram using 16 buckets
	compParams := hikaku.ComparisonParameters{BinsCount: 16, Threshold: 0.2}
	isImagesEqual, diff := hikaku.CompareByHistograms(first, second, compParams)

	// mix of previous two methods
	// 1. Checks width and height
	// 2. Compare by histograms
	isImagesEqual, diff := hikaku.Compare(first, second, compParams)
}
```

#### CompareByHistograms

**Description**
<br>Idea is simple:

1. Get histogram for Red, Green and Blue channels
2. Normalize histograms
3. Calculate [Hellinger distance](https://en.wikipedia.org/wiki/Hellinger_distance)
4. Summarize and make decision based on highlighted parts

### Difference highlighting

#### Using mask

Imagine you need to find a difference between two images:

<p align="center">
1.<img src="https://github.com/opa-oz/hikaku/raw/main/docs/golden.png" width="150" alt="Original"/>
2.<img src="https://github.com/opa-oz/hikaku/raw/main/docs/copper.png" width="150" alt="Changed"/>
</p>

Just run the code:

```go
package main

import (
	"github.com/opa-oz/hikaku"
	"image"
)

func getDifference(first, second image.Image) *image.Image {
	diffMask := hikaku.FindDiffMask(first, second, hikaku.ContextParameters{})
	straightForwardDiff := hikaku.ApplyDiff(second, diffMask, 128)

	return straightForwardDiff
}
```

**Voila!** Here is the difference, conveniently highlighted in red!
<p align="center">
<img src="https://github.com/opa-oz/hikaku/raw/main/docs/mask.png" width="150" alt="Changed"/>
<br>
</p>

#### Using shapes

Sometimes images with highlighted diff can be _noisy_. Let's fix it a little bit, combining pixels into bigger groups:

```go
func getDifference(first, second image.Image) *image.Image {
diffShapes := hikaku.FindDiffShapesMask(first, second,  hikaku.ContextParameters{})
shapesForwardDiff := hikaku.ApplyDiff(second, diffShapes, 128)

return straightForwardDiff
}
```

<p align="center">
<img src="https://github.com/opa-oz/hikaku/raw/main/docs/shape.png" width="150" alt="Changed"/><sub></sub>
</p>

One square to rule them all!

[**More examples**](https://github.com/opa-oz/hikaku/blob/main/demo/main.go)

### Optimization

If you need to compare huge amount of images, you may want to pre-calculate histograms.
```go
package main

import (
	"github.com/opa-oz/hikaku"
	"image"
)

func main(original, slightlyChanged, completelyDifferent image.Image) {
	compParams := hikaku.ComparisonParameters{BinsCount: 16}
	
	originalHist := hikaku.PrepareHistogram(original, compParams)
	slightlyChangedHist := hikaku.PrepareHistogram(slightlyChanged, compParams)
	completelyDifferentHist := hikaku.PrepareHistogram(completelyDifferent, compParams)

	// 1st vs 2nd
	isImagesEqual, diff := hikaku.CompareHistogramsOnly(originalHist, slightlyChangedHist, compParams)

	// 1st vs 3rd
	isImagesEqual, diff = hikaku.CompareHistogramsOnly(originalHist, completelyDifferentHist, compParams)

	// 2nd vs 3rd
	isImagesEqual, diff = hikaku.CompareHistogramsOnly(slightlyChangedHist, completelyDifferentHist, compParams)
	
	var unexpectedFourth image.Image

	// pass first image's histogram
	compParams.NormalizedGoldHist = originalHist
	
	// Yep, you still need to pass `original` image, but it will use pre-calculated histogram for this
	// while still calculate histogram for 4th image on demand
	isImagesEqual, diff = hikaku.CompareByHistograms(original, unexpectedFourth, compParams)
}
```

### Generated docs

```go
// Types
type Histogram = [][3]int
type NormalizedHistogram = [][3]float64

// Structs
type ComparisonParameters struct{ ... }
type ContextParameters struct{ ... }
type DiffPoint struct{ ... }
type Mask struct{ ... }
type Rect struct{ ... }

// Comparison
func Compare(golden, copper image.Image, params ComparisonParameters) (bool, float64)
func CompareByHistograms(golden, copper image.Image, params ComparisonParameters) (bool, float64)
func CompareHistogramsOnly(params ComparisonParameters) (bool, float64)
func CompareByParams(golden, copper image.Image) bool

// Overlay generation
func FindDiffMask(golden, copper image.Image, params ContextParameters) *image.RGBA
func FindDiffShapesMask(golden, copper image.Image, params ContextParameters) *image.RGBA
func ApplyDiff(golden, diff image.Image, transparency uint8) *image.RGBA

// Utils for optimization
func GetDiffShapes(golden, copper image.Image, params ContextParameters) map[int]Rect
func GetDiffPairs(golden, copper image.Image) []DiffPoint
func GetDiffPairsGrayscale(golden, copper image.Image) []DiffPoint
func PrepareHistogram(image image.Image, params ComparisonParameters) (normalizedHistogram NormalizedHistogram)
```

## Cheers 🥂

[*] _[Original Gophers' images repo](https://github.com/egonelbre/gophers)_