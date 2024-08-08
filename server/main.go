package main

import (
	"net/http"
	"os"

	"github.com/raphael-p/kafkito/server/config"
	"github.com/raphael-p/kafkito/server/resolvers"
	"github.com/raphael-p/kafkito/server/utils"
)

func main() {
	utils.InitLogger()
	defer utils.CloseLogger()
	mux := http.NewServeMux()

	mux.HandleFunc("POST /queue/{name}", resolvers.CreateQueue)
	mux.HandleFunc("POST /queue/{oldName}/rename/{newName}", resolvers.RenameQueue)
	mux.HandleFunc("DELETE /queue/{name}", resolvers.DeleteQueue)
	mux.HandleFunc("GET /queues", resolvers.ListQueues)
	mux.HandleFunc("POST /queue/{name}/publish", resolvers.PublishMessage)
	mux.HandleFunc("GET /queue/{name}/messages", resolvers.ReadMessages)
	mux.HandleFunc("DELETE /message/{id}", resolvers.ConsumeMessage)

	port := os.Getenv(config.PORT_ENVAR)
	if port == "" {
		port = config.DEFAULT_PORT
	}
	utils.LogTrace("server started on port " + port + "\n")
	err := http.ListenAndServe(":"+port, mux)
	utils.LogError(err.Error())
}
