package resolvers

import (
	"fmt"
	"os/exec"

	"github.com/raphael-p/kafkito/cli/utils"
)

func StartServer() {
	cmd := exec.Command("kafkitoserver")
	if err := cmd.Start(); err != nil {
		fmt.Println("error: could not execute kafkito's server binary:", err.Error())
		return
	}

	if err := utils.PingWithRetry(); err != nil {
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
