package utils

import (
	"fmt"
	"time"
)

func PingWithRetry() error {
	defaultError := fmt.Errorf("ping failed")
	for i := 1; i <= 5; i++ {
		time.Sleep(time.Millisecond * 200 * time.Duration(i)) // linear backoff
		response := KafkitoGet("/ping/kafkito")

		// retry until there is no error and up to five times
		if response.Error != nil && response.BodyString == "retry" && i < 5 {
			continue
		} else if response.Error != nil {
			defaultError = response.Error
			break
		} else {
			return nil
		}
	}
	return defaultError
}
