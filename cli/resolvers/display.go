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

func displayCSV(stream io.ReadCloser, columnWidth []int, formatter dataFormatter) {
	defer stream.Close()

	scanner := bufio.NewScanner(stream)
	headerRow := true
	for scanner.Scan() {
		row := scanner.Text()
		cells := strings.Split(row, ",")
		for columnIndex, cell := range cells {
			// print cell with custom formatting for non-header rows
			var formattedCell string
			if headerRow {
				headerRow = false
				formattedCell = cell
			} else {
				data, err := formatter(columnIndex, cell)
				if err != nil {
					fmt.Println(err)
					return
				}
				formattedCell = data
			}
			fmt.Print(formattedCell)

			// add consistent spacing to align the columns
			if columnIndex+1 < len(cells) {
				if columnIndex >= len(columnWidth) {
					break // new columns have been added unexpectedly
				}
				width := columnWidth[columnIndex]
				padding := float64(width - len(formattedCell))
				spaceCount := int(math.Max(0, padding)) + 2
				fmt.Print(strings.Repeat(" ", spaceCount))
			}
		}
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
