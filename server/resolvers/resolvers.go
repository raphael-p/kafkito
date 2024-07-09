package resolvers

import (
	"fmt"
	"net/http"

	"github.com/raphael-p/kafkito/server/config"
	"github.com/raphael-p/kafkito/server/queue"
)

func parseQueueName(w http.ResponseWriter, r *http.Request) (string, bool) {
	queueName := r.PathValue("name")
	errPrefix := "error parsing queue name: "

	if queueName == "" {
		errBody := errPrefix + "queue name must not be empty"
		http.Error(w, errBody, http.StatusBadRequest)
		return "", false
	}

	if len(queueName) > int(config.MAX_QUEUE_NAME_BYTES) {
		errBody := fmt.Sprint(
			errPrefix,
			"queue name is too long, max is: ",
			config.MAX_QUEUE_NAME_BYTES,
		)
		http.Error(w, errBody, http.StatusBadRequest)
		return "", false
	}

	return queueName, true
}

func CreateQueue(queues queue.QueueMap) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queueName, ok := parseQueueName(w, r)
		if !ok {
			return
		}

		err := queues.AddQueue(queueName)
		if err != nil {
			errBody := "error adding queue: " + err.Error()
			http.Error(w, errBody, http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func ListQueues(queues queue.QueueMap) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/csv")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("name,created_at\n"))
		for _, q := range queues {
			w.Write([]byte(fmt.Sprint(q.Name, ",", q.CreatedAt, "\n")))
		}
	}
}

func PublishMessage(queues queue.QueueMap) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queueName, ok := parseQueueName(w, r)
		if !ok {
			return
		}

		q, ok := queues[queueName]
		if !ok {
			errBody := "error fetching queue: no queue with name: " + queueName
			http.Error(w, errBody, http.StatusBadRequest)
			return
		}

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "error parsing request body", http.StatusBadRequest)
			return
		}
		header := r.FormValue("header")
		body := r.FormValue("body")

		message, err := queue.MakeMessage(header, body, config.MESSAGE_TTL)
		if err != nil {
			errBody := "error creating message: " + err.Error()
			http.Error(w, errBody, http.StatusBadRequest)
			return
		}

		q.Messages = append(q.Messages, message)
		queues[queueName] = q
		w.WriteHeader(http.StatusCreated)
	}
}
