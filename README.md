# Hikaku

<p align="center">Yet another tool for two image comparison.</p> 

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

### comparison

> noun (common) (futsuumeishi), noun or participle which takes the aux. verb suru, nouns which may take the genitive
> case particle `no'

## Motivation

**Hikaku** was made for my projects, focused on generating images. I **really** want to compare generated images or
screenshots **fast** and **reliable**.

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

**Description**: Idea is simple:

1. Get histogram for Red, Green and Blue channels
2. Normalize histograms
3. Calculate [Hellinger distance](https://en.wikipedia.org/wiki/Hellinger_distance)
4. Summarize and make decision


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