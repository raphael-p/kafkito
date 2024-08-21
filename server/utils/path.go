package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetDirectory() string {
	if os.Args[0] == "kafkitoserver" {
		// if executed from compiled binary, use its directory as the log location
		ex, err := os.Executable()
		if err != nil {
			panic(fmt.Sprintf("failed to locate executable: %s", err))
		}
		return filepath.Dir(ex)
	} else {
		// otherwise, use cwd as log location
		return "."
	}
}
