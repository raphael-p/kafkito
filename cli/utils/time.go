package utils

import "time"

const TIME_CHARS int = 26

func FormatTime(datetime time.Time) string {
	return datetime.Format("2006-01-02 15:04:05 -0700")
}
