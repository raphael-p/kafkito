package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func getConfigDirectory() (string, error) {
	if os.Args[0] == "kafkito" {
		// if executed from compiled binary, use its directory
		ex, err := os.Executable()
		if err != nil {
			return "", fmt.Errorf("error: failed to locate executable: %s", err)
		}
		return filepath.Dir(ex), nil
	} else {
		// otherwise, use parent directory
		return "..", nil
	}
}

type Config struct {
	Port         string `json:"port"`
	MaxQueueName uint8  `json:"max_queue_name_bytes"`
}

var config Config

func IntialiseConfig() error {
	configDir, err := getConfigDirectory()
	if err != nil {
		return err
	}
	configPath := filepath.Join(configDir, "config.json")
	file, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("error: failed to open config file: %s", err)
	}
	defer file.Close()

	var values *Config = &Config{}
	if err = json.NewDecoder(file).Decode(values); err != nil {
		return fmt.Errorf("error: could not parse config file: %s", err)
	}

	_, err = strconv.Atoi(values.Port)
	if err != nil {
		return fmt.Errorf("error: invalid port \"%s\"", values.Port)
	}

	if values.MaxQueueName <= 0 {
		return fmt.Errorf("error: max_queue_name_bytes must be specified in the config file and be greater than 0")
	}

	config = *values
	return nil
}

func GetPort() string {
	return config.Port
}

func GetQueueNameMaxLength() uint8 {
	return config.MaxQueueName
}
