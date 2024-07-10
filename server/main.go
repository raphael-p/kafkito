package main

import (
	"log"
	"net/http"
	"os"

	"github.com/raphael-p/kafkito/server/config"
	"github.com/raphael-p/kafkito/server/resolvers"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /createQueue/{name}", resolvers.CreateQueue)
	mux.HandleFunc("GET /queues", resolvers.ListQueues)
	mux.HandleFunc("POST /queue/{name}/publish", resolvers.PublishMessage)
	mux.HandleFunc("GET /queue/{name}/messages", resolvers.ReadMessages)

	port := os.Getenv(config.PORT_ENVAR)
	if port == "" {
		port = config.DEFAULT_PORT
	}
	log.Println("server started on port ", port)
	err := http.ListenAndServe(":"+port, mux)
	log.Fatal(err)
}
