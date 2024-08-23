package utils

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

var port string

func ValidatePort() bool {
	if port == "" {
		var err error
		port, err = readPortNumber()
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
	}
	return true
}

func GetPort() string {
	return port
}

type KafkitoResponse struct {
	StatusCode int
	Body       string
	Error      error
}

type MakeHTTPRequest func() (*http.Response, error)

func responseHandler(res *http.Response, callError error) KafkitoResponse {
	if callError != nil {
		return KafkitoResponse{0, "retry", callError}
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return KafkitoResponse{0, "", err}
	}

	return KafkitoResponse{res.StatusCode, string(body), nil}

}

func KafkitoGet(endpoint string) KafkitoResponse {
	return responseHandler(http.Get(
		"http://localhost:" + GetPort() + endpoint,
	))
}

func KafkitoPost(endpoint, reqContentType, reqBody string) KafkitoResponse {
	var reqBodyReader io.Reader = strings.NewReader(reqBody)
	return responseHandler(http.Post(
		"http://localhost:"+GetPort()+endpoint,
		reqContentType,
		reqBodyReader,
	))
}
