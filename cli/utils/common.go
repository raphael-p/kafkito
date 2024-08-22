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
		statusCode, body, err := KafkitoGet("/ping/kafkito")

		if err != nil && body == "retry" && i < 5 {
			// retry if there is an error
			continue
		} else if err != nil {
			// after fifth retry, show the error
			defaultError = err
			break
		} else if statusCode == http.StatusOK {
			// happy path
			return nil
		} else {
			// if server returns an error, show it
			return fmt.Errorf("status code %d: %s", statusCode, body)
		}
	}
	return defaultError
}
