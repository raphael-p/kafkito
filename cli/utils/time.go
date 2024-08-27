package utils

import (
	"fmt"
	"strconv"
	"time"
)

const TIME_FORMAT string = "2006-01-02 15:04:05 -0700"
const TIME_CHAR_COUNT int = len(TIME_FORMAT)

func UnixToDateTime(timestamp string) (string, error) {
	unixSeconds, err := strconv.Atoi(timestamp)
	if err != nil {
		return "", fmt.Errorf("error: %s", err)
	}
	datetime := time.Unix(int64(unixSeconds), 0)
	return datetime.Format(TIME_FORMAT), nil
}
