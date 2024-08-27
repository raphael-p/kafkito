package main

import (
	"context"
	"net/http"

	"github.com/raphael-p/kafkito/server/config"
	"github.com/raphael-p/kafkito/server/resolvers"
	"github.com/raphael-p/kafkito/server/utils"
)

var server *http.Server
var gracefulShutdown bool

func main() {
	utils.InitLogger()
	config.ReadConfigFile()
	defer utils.CloseLogger()
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping/kafkito", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("POST /shutdown", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		go func() {
			gracefulShutdown = true
			if err := server.Shutdown(context.Background()); err != nil {
				gracefulShutdown = false
				utils.LogError("server shutdown failed: " + err.Error())
			}
		}()
	})
	mux.HandleFunc("POST /queue/{name}", resolvers.CreateQueue)
	mux.HandleFunc("POST /queue/{oldName}/rename/{newName}", resolvers.RenameQueue)
	mux.HandleFunc("DELETE /queue/{name}", resolvers.DeleteQueue)
	mux.HandleFunc("GET /queues", resolvers.ListQueues)
	mux.HandleFunc("POST /queue/{name}/publish", resolvers.PublishMessage)
	mux.HandleFunc("GET /queue/{name}/messages", resolvers.ReadMessages)
	mux.HandleFunc("GET /message/{id}", resolvers.ReadMessage)
	mux.HandleFunc("DELETE /message/{id}", resolvers.ConsumeMessage)

	port := config.Values.Port
	server = &http.Server{Addr: ":" + port, Handler: mux}
	gracefulShutdown = false
	utils.LogTrace("server started on port " + port + "\n")
	err := server.ListenAndServe()

	if !gracefulShutdown || err.Error() != "http: Server closed" {
		utils.LogError(err.Error())
	} else {
		utils.LogTrace("server shutdown gracefully")
	}
}
