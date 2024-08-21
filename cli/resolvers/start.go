package resolvers

import (
	"fmt"
	"net/http"
	"os/exec"
	"time"

	"github.com/raphael-p/kafkito/cli/utils"
)

func StartServer() {
	cmd := exec.Command("kafkitoserver")
	if err := cmd.Start(); err != nil {
		fmt.Println("failed to start kafkito: ", err.Error())
		return
	}

	for i := 1; i <= 5; i++ {
		time.Sleep(time.Millisecond * 200 * time.Duration(i)) // linear backoff
		statusCode, body, err := utils.KafkitoGet("/ping/kafkito")

		if err != nil && body == "retry" && i < 5 {
			// retry if there is an error
			continue
		} else if err != nil {
			// after fifth retry, show the error
			fmt.Println("kafkito ping failed: ", err.Error())
		} else if statusCode == http.StatusOK {
			// happy path
			fmt.Println("kafkito is running")
		} else {
			// if server returns an error, show it
			fmt.Printf(
				"%d error occured while starting kafkito: %s\n",
				statusCode, body,
			)
		}
		break
	}
}
