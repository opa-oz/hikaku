package hikaku

import (
	"image"
)

type DiffPoint struct {
	x, y int
}

type Rect struct {
	minX, minY, maxX, maxY int
}

var dir = [...]struct{ first, second int }{
	{0, 1},
	{0, 2},
	{0, -1},
	{0, -2},
	{0, 0},
	{0, 1},
	{0, 2},
	{0, -1},
	{0, -2},
	{1, 2},
	{1, -1},
	{1, -2},
	{1, 0},
	{1, 1},
	{1, 2},
	{1, -1},
	{1, -2},
	{2, -1},
	{2, -2},
	{2, 0},
	{2, 1},
	{2, 2},
	{2, -1},
	{2, -2},
	{-1, -2},
	{-1, 0},
	{-1, 1},
	{-1, 2},
	{-1, -1},
	{-1, -2},
	{-2, 0},
	{-2, 1},
	{-2, 2},
	{-2, -1},
	{-2, -2},
	{0, 1},
	{0, 2},
	{0, -1},
	{0, -2},
	{1, 2},
	{1, -1},
	{1, -2},
	{2, -1},
	{2, -2},
	{-1, -2},
}

// dfs is an implementation of Depth-first search algorithm
// @see {@link https://en.wikipedia.org/wiki/Depth-first_search}
func dfs(matrix [][]bool, visited [][]int, x, y, n, m int, currentShape int) bool {
	if !check(x, y, n, m) || !matrix[x][y] || visited[x][y] != 0 {
		return false
	}

	visited[x][y] = currentShape

	for _, element := range dir {
		x1 := x + element.first
		y1 := y + element.second

		dfs(matrix, visited, x1, y1, n, m, currentShape)
	}

	return true
}

// GetDiffPairs gets pairs of pixels, that are different in the both images
// It returns slice of pairs (x,y), representing pixel's coordinates
// Works with both Gray and RGBA images
func GetDiffPairs(golden, copper image.Image) []DiffPoint {
	bounds := golden.Bounds()
	result := make([]DiffPoint, bounds.Max.Y*bounds.Max.X)

	amount := 0

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if golden.At(x, y) != copper.At(x, y) {
				result[amount] = DiffPoint{x, y}
				amount += 1
			}
		}
	}

	return result[0:amount]
}

// GetDiffShapes finds shapes (rectangles) inside diff to avoid "too noisy" diff as a result
func GetDiffShapes(golden, copper image.Image, params ContextParameters) map[int]Rect {
	bounds := golden.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	var diffPairs []DiffPoint

	if params.DiffPairs != nil {
		diffPairs = params.DiffPairs
	} else {
		diffPairs = GetDiffPairsGrayscale(golden, copper)
	}

	matrix := make([][]bool, width)
	visited := make([][]int, width)
	for i := range matrix {
		matrix[i] = make([]bool, height)
		visited[i] = make([]int, height)
	}

	for i := range diffPairs {
		pair := diffPairs[i]
		matrix[pair.x][pair.y] = true
	}

	currentShape := 1

	for _, pair := range diffPairs {
		found := dfs(matrix, visited, pair.x, pair.y, width, height, currentShape)

		if found {
			currentShape++
		}
	}

	groups := make(map[int]Rect)

	for _, pair := range diffPairs {
		index := visited[pair.x][pair.y]
		value, ok := groups[index]

		if !ok {
			groups[index] = Rect{pair.x, pair.y, pair.x, pair.y}
		} else {
			value.minX = min(value.minX, pair.x)
			value.maxX = max(value.maxX, pair.x)
			value.minY = min(value.minY, pair.y)
			value.maxY = max(value.maxY, pair.y)

			groups[index] = value
		}
	}

	return groups
}
