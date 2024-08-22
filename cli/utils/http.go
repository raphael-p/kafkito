package utils

import (
	"io"
	"net/http"
)

var port string

func getPort() (string, error) {
	if port == "" {
		var err error
		port, err = readPortNumber()
		if err != nil {
			return "", err
		}
	}
	return port, nil
}

type KafkitoResponse struct {
	StatusCode int
	Body       string
	Error      error
}

func KafkitoGet(endpoint string) KafkitoResponse {
	port, err := getPort()
	if err != nil {
		return KafkitoResponse{0, "", err}
	}

	res, err := http.Get("http://localhost:" + port + endpoint)
	if err != nil {
		return KafkitoResponse{0, "retry", err}
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return KafkitoResponse{0, "", err}
	}

	return KafkitoResponse{res.StatusCode, string(body), nil}
}
