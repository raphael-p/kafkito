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

func KafkitoGet(endpoint string) (int, string, error) {
	port, err := getPort()
	if err != nil {
		return 0, "", err
	}

	res, err := http.Get("http://localhost:" + port + endpoint)
	if err != nil {
		return 0, "retry", err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, "", err
	}

	return res.StatusCode, string(body), nil
}
