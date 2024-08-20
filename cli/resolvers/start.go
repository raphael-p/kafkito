package resolvers

import (
	"os/exec"
)

func StartServer() {
	exec.Command("kafkitoserver", "&")
	// TODO: run ping for confirmation
	// fmt.Println("server started")
}
