package resolvers

import (
	"bufio"
	"fmt"
	"io"
	"math"
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

type dataFormatter func(index int, data string) (string, error)

func displayCSV(stream io.ReadCloser, columnWidth int, formatter dataFormatter) {
	defer stream.Close()

	scanner := bufio.NewScanner(stream)
	headerRow := true
	for scanner.Scan() {
		row := scanner.Text()
		cells := strings.Split(row, ",")
		for index, cell := range cells {
			if headerRow {
				fmt.Print(cell)
			} else {
				data, err := formatter(index, cell)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Print(data)
			}

			if index+1 < len(cells) {
				spaceCount := int(math.Max(0, float64(columnWidth-len(cell)))) + 2
				fmt.Print(strings.Repeat(" ", spaceCount))
			}
		}
		headerRow = false
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
	paginate(help)
}
