package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type KafkitoResponse struct {
	StatusCode int
	BodyStream io.ReadCloser
	BodyString string
	Header     http.Header
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
		var errStr string
		errBody, err := io.ReadAll(res.Body)

		if err == nil {
			errStr = string(errBody)
		} else {
			errStr = err.Error()
		}
		errStr = strings.TrimSpace(errStr)

		return "", fmt.Errorf(
			"error: status code %d: %s",
			res.StatusCode, errStr,
		)
	}

	return "", nil
}

func responseHandlerString(res *http.Response, callError error) KafkitoResponse {
	if body, err := responseErrorHandler(res, callError); err != nil {
		return KafkitoResponse{0, nil, body, nil, err}
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return KafkitoResponse{0, nil, "", nil, fmt.Errorf("error: %s", err)}
	}

	return KafkitoResponse{res.StatusCode, nil, string(body), res.Header, nil}
}

func responseHandlerStream(res *http.Response, callError error) KafkitoResponse {
	if body, err := responseErrorHandler(res, callError); err != nil {
		return KafkitoResponse{0, nil, body, nil, err}
	}

	return KafkitoResponse{res.StatusCode, res.Body, "", res.Header, nil}
}

func KafkitoGet(endpoint string) KafkitoResponse {
	return responseHandlerString(http.Get(
		"http://localhost:" + GetPort() + endpoint,
	))
}

func KafkitoGetStream(endpoint string) KafkitoResponse {
	return responseHandlerStream(http.Get(
		"http://localhost:" + GetPort() + endpoint,
	))
}

func KafkitoPost(endpoint string) KafkitoResponse {
	return responseHandlerString(http.Post(
		"http://localhost:"+GetPort()+endpoint,
		"",
		nil,
	))
}

func KafkitoPostForm(endpoint string, data url.Values) KafkitoResponse {
	return responseHandlerString(http.PostForm(
		"http://localhost:"+GetPort()+endpoint,
		data,
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

func KakitoDeleteStream(endpoint string) KafkitoResponse {
	client := &http.Client{}
	req, _ := http.NewRequest(
		http.MethodDelete,
		"http://localhost:"+GetPort()+endpoint,
		nil,
	)
	return responseHandlerStream(client.Do(req))
}

func IsSuccessful(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}
