package utils

import (
	"io"
	"net/http"
	"strings"
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

type MakeHTTPRequest func(port string) (*http.Response, error)

func kafkitoHTTP(makeRequest MakeHTTPRequest) KafkitoResponse {
	port, err := getPort()
	if err != nil {
		return KafkitoResponse{0, "", err}
	}

	res, err := makeRequest(port)
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

func KafkitoGet(endpoint string) KafkitoResponse {
	return kafkitoHTTP(func(port string) (*http.Response, error) {
		return http.Get("http://localhost:" + port + endpoint)
	})
}

func KafkitoPost(endpoint, reqContentType, reqBody string) KafkitoResponse {
	return kafkitoHTTP(func(port string) (*http.Response, error) {
		var reqBodyReader io.Reader = strings.NewReader(reqBody)
		return http.Post(
			"http://localhost:"+port+endpoint,
			reqContentType,
			reqBodyReader,
		)
	})
}
