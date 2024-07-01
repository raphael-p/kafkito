package main

import (
	"log"
	"net/http"
	"os"

	"github.com/raphael-p/kafkito/server/config"
	"github.com/raphael-p/kafkito/server/queue"
	"github.com/raphael-p/kafkito/server/resolvers"
)

func main() {
	queueMap := make(queue.QueueMap)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /createQueue/{name}", resolvers.CreateQueue(queueMap))

	port := os.Getenv(config.PORT_ENVAR)
	if port == "" {
		port = config.DEFAULT_PORT
	}
	log.Println("server started on port ", port)
	err := http.ListenAndServe(":"+port, mux)
	log.Fatal(err)
}
