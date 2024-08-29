package resolvers

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func usePager(text string) {
	pager := exec.Command("less")
	pager.Stdin = strings.NewReader(text)
	pager.Stdout = os.Stdout
	pager.Stderr = os.Stderr
	if err := pager.Run(); err != nil {
		fmt.Print(text)
	}
}

type PrintRow func(row string, isHeader bool) bool

func displayCSV(stream io.ReadCloser, printRow PrintRow, skipHeader bool) {
	defer stream.Close()

	scanner := bufio.NewScanner(stream)
	isHeader := true
	for scanner.Scan() {
		if skipHeader && isHeader {
			isHeader = false
			continue
		}

		row := scanner.Text()
		if !printRow(row, isHeader) {
			return
		}
		isHeader = false
		fmt.Print("\n")
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error: could not read CSV response:", err.Error())
	}
}

func DisplaySeekHelp(prefix string) {
	var space string
	if prefix != "" {
		space = " "
	}
	fmt.Print(prefix, space, "To see all valid commands, run 'kafkito help'\n")
}

func DisplayHelp() {
	usePager(help)
}
