package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/raphael-p/kafkito/server/utils"
)

var Values *config

type config struct {
	Port                  string `json:"port"`
	MaxQueues             uint32 `json:"max_queues"`
	MaxQueueLength        uint32 `json:"max_queue_length"`
	MaxQueueNameBytes     uint32 `json:"max_queue_name_bytes"`
	MaxMessageHeaderBytes uint32 `json:"max_message_header_bytes"`
	MaxMessageBodyBytes   uint32 `json:"max_message_body_bytes"`
	MessageTTL            int64  `json:"message_ttl"`
	MessageBatchSize      uint32 `json:"message_batch_size"`
}

func ReadConfigFile() {
	logFilepath := filepath.Join(utils.GetDirectory(), "config.json")
	file, err := os.Open(logFilepath)
	if err != nil {
		panic(fmt.Sprintf("failed to open config files: %s", err))
	}
	defer file.Close()

	// Unmarshal the JSON into the struct
	Values = &config{}
	if err = json.NewDecoder(file).Decode(Values); err != nil {
		panic(fmt.Sprint("could not parse config file: ", err))
	}
}
