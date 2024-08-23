package resolvers

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/raphael-p/kafkito/cli/utils"
)

func StartServer() {
	cmd := exec.Command("kafkitoserver")
	if err := cmd.Start(); err != nil {
		fmt.Println("could not execute kafkito's server binary:", err.Error())
		return
	}

	if err := utils.PingWithRetry(); err != nil {
		fmt.Println("could not ping kafkito:", err.Error())
		return
	}

	fmt.Println("kafkito is runing on port", utils.GetPort())
}

func StopServer() {
	response := utils.KafkitoPost("/shutdown", "", "")

	if response.Error != nil {
		fmt.Println(
			"kafkito not available (maybe already stopped):",
			response.Error.Error(),
		)
		return
	}

	if response.StatusCode != http.StatusAccepted {
		fmt.Printf(
			"kafkito could not be stopped: status code %d: %s\n",
			response.StatusCode, response.Body,
		)
		return
	}

	fmt.Println("kafkito was stopped")
}

func ServerInfo() {
	response := utils.KafkitoGet("/ping/kafkito")

	if response.Error != nil {
		fmt.Println("kafkito is not running on port", utils.GetPort())
		return
	}

	if response.StatusCode != http.StatusOK {
		fmt.Printf(
			"status code %d: %s\n",
			response.StatusCode, response.Body,
		)
		return
	}

	fmt.Println("kafkito is runing on port", utils.GetPort())
}
