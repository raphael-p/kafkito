package resolvers

import (
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"time"
)

func ping() (int, string, error) {
	// TODO: un-hardcode port number
	res, err := http.Get("http://localhost:8083/ping/kafkito")
	if err != nil {
		return 0, "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, "", err
	}
	return res.StatusCode, string(body), nil
}

func StartServer() {
	cmd := exec.Command("kafkitoserver")
	if err := cmd.Start(); err != nil {
		fmt.Println("failed to start kafkito: ", err.Error())
	} else {
		for i := 1; i <= 5; i++ {
			time.Sleep(time.Millisecond * 200 * time.Duration(i)) // linear backoff
			statusCode, body, err := ping()
			if err != nil && i < 5 {
				continue
			} else if err != nil {
				fmt.Println("kafkito ping failed after 3s: ", err.Error())
			} else if statusCode == http.StatusOK {
				fmt.Println("kafkito is running")
			} else {
				fmt.Printf(
					"%d error occured while starting kafkito: %s\n",
					statusCode, body,
				)
			}
			break
		}
	}
}
