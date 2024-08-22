package main

import (
	"flag"

	"github.com/raphael-p/kafkito/cli/resolvers"
)

const HELP = "help"
const START_SERVER = "start"
const STOP_SERVER = "stop"
const CREATE_QUEUE = "create"
const RENAME_QUEUE = "rename"
const DELETE_QUEUE = "delete"
const LIST_QUEUES = "list"
const PUBLISH_MESSAGE = "publish"
const READ_QUEUE = "read"
const CONSUME_MESSAGE = "consume"

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		resolvers.DisplaySeekHelp("Welcome to Kafkito!")
	} else if flag.Arg(0) == START_SERVER {
		resolvers.StartServer()
	} else if flag.Arg(0) == STOP_SERVER {
		resolvers.StopServer()
	} else if flag.Arg(0) == HELP {
		resolvers.DisplayHelp()
	} else if flag.Arg(0) == CREATE_QUEUE {
		resolvers.CreateQueue()
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
