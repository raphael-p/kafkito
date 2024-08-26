package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type KafkitoResponse struct {
	StatusCode int
	BodyStream io.ReadCloser
	BodyString string
	Error      error
}

type MakeHTTPRequest func() (*http.Response, error)

func responseErrorHandler(res *http.Response, callError error) (string, error) {
	if callError != nil {
		return "retry", fmt.Errorf(
			"error: kafkito is not running on port %s",
			GetPort(),
		)
	}

	if !IsSuccessful(res.StatusCode) {
		defer res.Body.Close()
		errBody, err := io.ReadAll(res.Body)

		if err == nil {
			err = errors.New(string(errBody))
		}

		return "", fmt.Errorf(
			"error: status code %d: %s",
			res.StatusCode, err,
		)
	}

	return "", nil
}

func responseHandlerString(res *http.Response, callError error) KafkitoResponse {
	if body, err := responseErrorHandler(res, callError); err != nil {
		return KafkitoResponse{0, nil, body, err}
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return KafkitoResponse{0, nil, "", fmt.Errorf("error: %s", err)}
	}

	return KafkitoResponse{res.StatusCode, nil, string(body), nil}
}

func responseHandlerStream(res *http.Response, callError error) KafkitoResponse {
	if body, err := responseErrorHandler(res, callError); err != nil {
		return KafkitoResponse{0, nil, body, err}
	}

	return KafkitoResponse{res.StatusCode, res.Body, "", nil}
}

func KafkitoGet(endpoint string) KafkitoResponse {
	return responseHandlerString(http.Get(
		"http://localhost:" + GetPort() + endpoint,
	))
}

func KafkitoGetCSV(endpoint string) KafkitoResponse {
	return responseHandlerStream(http.Get(
		"http://localhost:" + GetPort() + endpoint,
	))
}

func KafkitoPost(endpoint, reqContentType, reqBody string) KafkitoResponse {
	var reqBodyReader io.Reader = strings.NewReader(reqBody)
	return responseHandlerString(http.Post(
		"http://localhost:"+GetPort()+endpoint,
		reqContentType,
		reqBodyReader,
	))
}

func KakitoDelete(endpoint string) KafkitoResponse {
	client := &http.Client{}
	req, _ := http.NewRequest(
		http.MethodDelete,
		"http://localhost:"+GetPort()+endpoint,
		nil,
	)
	return responseHandlerString(client.Do(req))
}

func IsSuccessful(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}
