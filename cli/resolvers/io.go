package resolvers

import (
	"fmt"
	"os/exec"

	"github.com/raphael-p/kafkito/cli/utils"
)

func StartServer() {
	cmd := exec.Command("kafkitoserver")
	if err := cmd.Start(); err != nil {
		fmt.Println("could not execute kafkito server binary: ", err.Error())
		return
	}

	if err := utils.PingWithRetry(); err != nil {
		fmt.Println("could not ping kafkito server: ", err.Error())
		return
	}

	fmt.Println("kafkito is running")
}
