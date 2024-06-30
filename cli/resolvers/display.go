package resolvers

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func paginate(text string) {
	pager := exec.Command("less")
	pager.Stdin = strings.NewReader(text)
	pager.Stdout = os.Stdout
	pager.Stderr = os.Stderr
	if err := pager.Run(); err != nil {
		fmt.Print(text)
	}
}

func DiplaySeekHelp(prefix string) {
	var space string
	if prefix != "" {
		space = " "
	}
	fmt.Print(prefix, space, "To see all valid commands, run 'kafkito help'\n")
}

func DisplayHelp() {
	paginate(help)
}

func DisplayCreate() {
	fmt.Print("placeholder for 'create' command\n")
}

func DisplayDelete() {
	fmt.Print("placeholder for 'delete' command\n")
}

func DisplayList() {
	fmt.Print("placeholder for 'list' command\n")
}

func DisplayPublish() {
	fmt.Print("placeholder for 'publish' command\n")
}

func DisplaySubscribe() {
	fmt.Print("placeholder for 'subscribe' command\n")
}
