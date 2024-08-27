package utils

const DEFAULT_WIDTH = 15
const MAX_BODY_DISPLAY = 30

func CalculateWidth(columnName string, expectedWidth int) int {
	if expectedWidth == -1 {
		return DEFAULT_WIDTH
	}

	if len(columnName) > expectedWidth {
		return len(columnName)
	}

	return expectedWidth
}
