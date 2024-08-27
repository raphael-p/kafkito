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

const DEFAULT_WIDTH = 15

func displayCSV(stream io.ReadCloser, columnWidth []int, formatter dataFormatter) {
	defer stream.Close()

	scanner := bufio.NewScanner(stream)
	headerRow := true
	for scanner.Scan() {
		row := scanner.Text()
		cells := strings.Split(row, ",")
		for colIdx, cell := range cells {
			// print cell with custom formatting for non-header rows
			if headerRow {
				fmt.Print(cell)
			} else {
				data, err := formatter(colIdx, cell)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Print(data)
			}

			// add consistent spacing to align the columns
			if colIdx+1 < len(cells) {
				var width int
				if colIdx < len(columnWidth) {
					width = columnWidth[colIdx]
				} else {
					width = DEFAULT_WIDTH
				}

				spaceCount := int(math.Max(0, float64(width-len(cell)))) + 2
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
