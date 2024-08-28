package main

import (
	"flag"
	"fmt"
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
	flag.Parse()
	command := flag.Arg(0)

	// noop commands
	if flag.NArg() == 0 {
		resolvers.DisplaySeekHelp("Welcome to Kafkito!")
		return
	} else if command == HELP {
		resolvers.DisplayHelp()
		return
	}

	if err := utils.IntialiseConfig(); err != nil {
		fmt.Println(err.Error())
		return
	}

	// all other commands (they require a valid port)
	switch command {
	case START_SERVER:
		resolvers.StartServer()
	case STOP_SERVER:
		resolvers.StopServer()
	case SERVER_INFO:
		resolvers.ServerInfo()
	case CREATE_QUEUE:
		if !validateArgs("queueName") {
			fmt.Println("usage: kafkito create <queueName>")
			return
		}
		resolvers.CreateQueue(flag.Arg(1))
	case RENAME_QUEUE:
		if !validateArgs("oldQueueName", "newQueueName") {
			fmt.Println("usage: kafkito rename <oldQueueName> <newQueueName>")
			return
		}
		resolvers.RenameQueue(flag.Arg(1), flag.Arg(2))
	case DELETE_QUEUE:
		if !validateArgs("queueName") {
			fmt.Println("usage: kafkito delete <queueName>")
			return
		}
		resolvers.DeleteQueue(flag.Arg(1))
	case LIST:
		if flag.Arg(1) == "" {
			resolvers.ListQueues() // list queues
		} else {
			resolvers.ReadMessages(flag.Arg(1)) // list messages of queue
		}
	case READ_MESSAGE:
		if !validateArgs("messageID") {
			fmt.Println("usage: kafkito read <messageID>")
			return
		}
		resolvers.ReadMessage(flag.Arg(1))
	case PUBLISH_MESSAGE:
		if !validateArgs("queueName", "message_header", "message_body") {
			fmt.Println("usage: kafkito publish <queueName> <message_header> <message_body>")
			return
		}
		message := strings.Join(flag.Args()[3:], " ")
		resolvers.PublishMessage(flag.Arg(1), flag.Arg(2), message)
	case CONSUME_MESSAGE:
		if !validateArgs("messageID") {
			fmt.Println("usage: kafkito consume <messageID>")
			return
		}
		resolvers.ConsumeMessage(flag.Arg(1))
	default:
		resolvers.DisplaySeekHelp("Command not recognised.")
	}
}

func validateArgs(args ...string) bool {
	for idx := range args {
		if flag.Arg(idx+1) == "" {
			fmt.Println(
				"missing arg(s):",
				strings.Join(args[idx:], ", "),
			)
			return false
		}
	}
	return true
}
