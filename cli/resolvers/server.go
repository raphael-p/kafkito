package resolvers

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/raphael-p/kafkito/cli/utils"
)

func pingWithRetry() error {
	defaultError := fmt.Errorf("ping failed")
	for i := 1; i <= 5; i++ {
		time.Sleep(time.Millisecond * 200 * time.Duration(i)) // linear backoff
		response := utils.KafkitoGet("/ping/kafkito")

		// retry until there is no error and up to five times
		if response.Error != nil && response.BodyString == "retry" && i < 5 {
			continue
		} else if response.Error != nil {
			defaultError = response.Error
			break
		} else {
			return nil
		}
	}
	return defaultError
}

func StartServer() {
	cmd := exec.Command("kafkitoserver")
	if err := cmd.Start(); err != nil {
		fmt.Println("error: could not execute kafkito's server binary:", err.Error())
		return
	}

	if err := pingWithRetry(); err != nil {
		fmt.Println("could not ping kafkito:", err.Error())
		return
	}

	fmt.Println("kafkito is running on port", utils.GetPort())
}

func StopServer() {
	response := utils.KafkitoPost("/shutdown")
	if response.Error != nil {
		fmt.Println(response.Error.Error())
		return
	}

	fmt.Println("kafkito was stopped")
}

func ServerInfo() {
	response := utils.KafkitoGet("/ping/kafkito")
	if response.Error != nil {
		fmt.Println(response.Error.Error())
		return
	}

	fmt.Println("kafkito is running on port", utils.GetPort())
}
