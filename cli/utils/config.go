package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type config struct {
	Port string `json:"port"`
}

func getConfigDirectory() (string, error) {
	if os.Args[0] == "kafkito" {
		// if executed from compiled binary, use its directory
		ex, err := os.Executable()
		if err != nil {
			return "", fmt.Errorf("failed to locate executable: %s", err)
		}
		return filepath.Dir(ex), nil
	} else {
		// otherwise, use parent directory
		return "..", nil
	}
}

func readPortNumber() (string, error) {
	configDir, err := getConfigDirectory()
	if err != nil {
		return "", err
	}
	configPath := filepath.Join(configDir, "config.json")
	file, err := os.Open(configPath)
	if err != nil {
		return "", fmt.Errorf("failed to open config file: %s", err)
	}
	defer file.Close()

	var values *config = &config{}
	if err = json.NewDecoder(file).Decode(values); err != nil {
		return "", fmt.Errorf("could not parse config file: %s", err)
	}

	_, err = strconv.Atoi(values.Port)
	if err != nil {
		return "", fmt.Errorf("invalid port \"%s\"", values.Port)
	}
	return values.Port, nil
}
