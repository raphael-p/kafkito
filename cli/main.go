package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/raphael-p/kafkito/cli/resolvers"
	"github.com/raphael-p/kafkito/cli/utils"
)

const HELP = "help"
const START_SERVER = "start"
const STOP_SERVER = "stop"
const SERVER_INFO = "info"
const CREATE_QUEUE = "create"
const RENAME_QUEUE = "rename"
const DELETE_QUEUE = "delete"
const LIST = "list"
const PUBLISH_MESSAGE = "publish"
const READ_MESSAGE = "read"
const CONSUME_MESSAGE = "consume"

func main() {
	if len(os.Args) < 2 {
		resolvers.DisplaySeekHelp("Welcome to Kafkito!")
		return
	}

	if err := utils.IntialiseConfig(); err != nil {
		fmt.Println(err.Error())
		return
	}

	command := os.Args[1]
	fs := flag.NewFlagSet(command, flag.ContinueOnError)
	fs.Parse(os.Args[2:])
	switch command {
	case HELP:
		resolvers.DisplayHelp()
	case START_SERVER:
		resolvers.StartServer()
	case STOP_SERVER:
		resolvers.StopServer()
	case SERVER_INFO:
		resolvers.ServerInfo()
	case CREATE_QUEUE:
		if !validateArgs(fs, "queueName") {
			fmt.Println("usage: kafkito create <queueName>")
			return
		}
		resolvers.CreateQueue(fs.Arg(0))
	case RENAME_QUEUE:
		if !validateArgs(fs, "oldQueueName", "newQueueName") {
			fmt.Println("usage: kafkito rename <oldQueueName> <newQueueName>")
			return
		}
		resolvers.RenameQueue(fs.Arg(0), fs.Arg(1))
	case DELETE_QUEUE:
		if !validateArgs(fs, "queueName") {
			fmt.Println("usage: kafkito delete <queueName>")
			return
		}
		resolvers.DeleteQueue(fs.Arg(0))
	case LIST:
		if flag.Arg(1) == "" {
			resolvers.ListQueues() // list queues
		} else {
			resolvers.ReadMessages(fs.Arg(0)) // list messages of queue
		}
	case READ_MESSAGE:
		if !validateArgs(fs, "messageID") {
			fmt.Println("usage: kafkito read <messageID>")
			return
		}
		resolvers.ReadMessage(fs.Arg(0))
	case PUBLISH_MESSAGE:
		if !validateArgs(fs, "queueName", "message_header", "message_body") {
			fmt.Println("usage: kafkito publish <queueName> <message_header> <message_body>")
			return
		}
		message := strings.Join(fs.Args()[2:], " ")
		resolvers.PublishMessage(fs.Arg(0), fs.Arg(1), message)
	case CONSUME_MESSAGE:
		if !validateArgs(fs, "messageID") {
			fmt.Println("usage: kafkito consume <messageID>")
			return
		}
		resolvers.ConsumeMessage(fs.Arg(0))
	default:
		resolvers.DisplaySeekHelp("Command not recognised.")
	}
}

func validateArgs(fs *flag.FlagSet, args ...string) bool {
	for idx := range args {
		if fs.Arg(idx) == "" {
			fmt.Println(
				"missing arg(s):",
				strings.Join(args[idx:], ", "),
			)
			return false
		}
	}
	return true
}
