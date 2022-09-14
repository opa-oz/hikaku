package hikaku

func check(x, y, n, m int) bool {
	return x >= 0 && y >= 0 && x < n-1 && y < m-1
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func max(a, b int) int {
	if a < b {
		return b
	}

	return a
}
