package utils

func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}

	if maxLength <= 3 {
		return "..."
	}

	return s[:maxLength-3] + "..."
}
