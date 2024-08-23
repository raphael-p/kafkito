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
const LIST_QUEUES = "list"
const PUBLISH_MESSAGE = "publish"
const READ_QUEUE = "read"
const CONSUME_MESSAGE = "consume"

func main() {
	flag.Parse()

	// noop commands
	if flag.NArg() == 0 {
		resolvers.DisplaySeekHelp("Welcome to Kafkito!")
		return
	} else if flag.Arg(0) == HELP {
		resolvers.DisplayHelp()
		return
	}

	if !utils.ValidatePort() {
		return
	}

	// all other commands (they require a valid port)
	if flag.Arg(0) == START_SERVER {
		resolvers.StartServer()
	} else if flag.Arg(0) == STOP_SERVER {
		resolvers.StopServer()
	} else if flag.Arg(0) == SERVER_INFO {
		resolvers.ServerInfo()
	} else if flag.Arg(0) == HELP {
		resolvers.DisplayHelp()
	} else if flag.Arg(0) == CREATE_QUEUE {
		if !validateArgs("queueName") {
			fmt.Println("expected: kafkito create <queueName>")
			return
		}
		resolvers.CreateQueue(flag.Arg(1))
	} else if flag.Arg(0) == DELETE_QUEUE {
		resolvers.DeleteQueue()
	} else if flag.Arg(0) == LIST_QUEUES {
		resolvers.ListQueues()
	} else if flag.Arg(0) == PUBLISH_MESSAGE {
		resolvers.PublishQueue()
	} else if flag.Arg(0) == READ_QUEUE {
		resolvers.ReadMessages()
	} else if flag.Arg(0) == CONSUME_MESSAGE {
		resolvers.ConsumeMessage()
	} else {
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
