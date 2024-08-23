package utils

import (
	"fmt"
	"net/http"
	"time"
)

func PingWithRetry() error {
	defaultError := fmt.Errorf("ping failed")
	for i := 1; i <= 5; i++ {
		time.Sleep(time.Millisecond * 200 * time.Duration(i)) // linear backoff
		response := KafkitoGet("/ping/kafkito")

		if response.Error != nil && response.Body == "retry" && i < 5 {
			// retry if there is an error
			continue
		} else if response.Error != nil {
			// after fifth retry, show the error
			defaultError = response.Error
			break
		} else if response.StatusCode == http.StatusOK {
			// happy path
			return nil
		} else {
			// if server returns an error, show it
			return fmt.Errorf(
				"status code %d: %s",
				response.StatusCode, response.Body,
			)
		}
	}
	return defaultError
}
