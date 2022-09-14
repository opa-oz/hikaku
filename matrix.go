package hikaku

import (
	"image"
)

type DiffPoint struct {
	X, Y int
}

type Rect struct {
	MinX, MinY, MaxX, MaxY int
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
		matrix[pair.X][pair.Y] = true
	}

	currentShape := 1

	for _, pair := range diffPairs {
		found := dfs(matrix, visited, pair.X, pair.Y, width, height, currentShape)

		if found {
			currentShape++
		}
	}

	groups := make(map[int]Rect)

	for _, pair := range diffPairs {
		index := visited[pair.X][pair.Y]
		value, ok := groups[index]

		if !ok {
			groups[index] = Rect{pair.X, pair.Y, pair.X, pair.Y}
		} else {
			value.MinX = min(value.MinX, pair.X)
			value.MaxX = max(value.MaxX, pair.X)
			value.MinY = min(value.MinY, pair.Y)
			value.MaxY = max(value.MaxY, pair.Y)

			groups[index] = value
		}
	}

	return groups
}
