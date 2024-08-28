package utils

import (
	"fmt"
	"math"
	"strings"
	"time"
)

const DEFAULT_WIDTH = 15
const MAX_BODY_DISPLAY = 30
const TIME_FORMAT string = "2006-01-02 15:04:05 -0700"
const TIME_CHAR_COUNT int = len(TIME_FORMAT)

func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}

	if maxLength <= 3 {
		return "..."
	}

	return s[:maxLength-3] + "..."
}

func UnixToDateTime(timestamp int) string {
	datetime := time.Unix(int64(timestamp), 0)
	return datetime.Format(TIME_FORMAT)
}

func CalculateWidth(columnName string, expectedWidth int) int {
	if expectedWidth == -1 {
		return DEFAULT_WIDTH
	}

	if len(columnName) > expectedWidth {
		return len(columnName)
	}

	return expectedWidth
}

func PrintCell(cell string, width int) {
	padding := float64(width - len(cell))
	spaceCount := int(math.Max(0, padding)) + 2
	trimmedCell := truncateString(cell, width)
	fmt.Printf("%s%s", trimmedCell, strings.Repeat(" ", spaceCount))
}
