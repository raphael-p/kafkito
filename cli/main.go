package main

import (
	"flag"

	"github.com/raphael-p/kafkito/cli/resolvers"
)

const HELP = "help"
const CREATE = "create"
const DELETE = "delete"
const LIST = "list"
const PUBLISH = "publish"
const SUBSCRIBE = "subscribe"

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		resolvers.DiplaySeekHelp("Welcome to Kafkito!")
	} else if flag.Arg(0) == HELP {
		resolvers.DisplayHelp()
	} else if flag.Arg(0) == CREATE {
		resolvers.DisplayCreate()
	} else if flag.Arg(0) == DELETE {
		resolvers.DisplayDelete()
	} else if flag.Arg(0) == LIST {
		resolvers.DisplayList()
	} else if flag.Arg(0) == PUBLISH {
		resolvers.DisplayPublish()
	} else if flag.Arg(0) == SUBSCRIBE {
		resolvers.DisplaySubscribe()
	} else {
		resolvers.DiplaySeekHelp("Command not recognised.")
	}
}
