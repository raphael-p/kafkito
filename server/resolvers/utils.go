package resolvers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/raphael-p/kafkito/server/config"
	"github.com/raphael-p/kafkito/server/queue"
)

func parseQueueName(w http.ResponseWriter, r *http.Request, pattern string) (string, bool) {
	queueName := r.PathValue(pattern)
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

func parseMessageID(w http.ResponseWriter, r *http.Request) (uint64, bool) {
	messageIDStr := r.PathValue("id")
	errPrefix := "error parsing message ID: "

	if messageIDStr == "" {
		errBody := errPrefix + "must not be empty"
		http.Error(w, errBody, http.StatusBadRequest)
		return 0, false
	}

	messageIDInt, err := strconv.Atoi(messageIDStr)
	if err != nil {
		errBody := errPrefix + err.Error()
		http.Error(w, errBody, http.StatusBadRequest)
		return 0, false
	}
	if messageIDInt <= 0 {
		errBody := errPrefix + "must be greater than zero"
		http.Error(w, errBody, http.StatusBadRequest)
		return 0, false
	}

	return uint64(messageIDInt), true
}

func getQueue(w http.ResponseWriter, r *http.Request, queues queue.QueueMap, pattern string) (queue.Queue, bool) {
	var q queue.Queue

	queueName, ok := parseQueueName(w, r, pattern)
	if !ok {
		return q, false
	}

	q, ok = queues.GetQueue(queueName)
	if !ok {
		errBody := "error fetching queue: no queue with name: " + queueName
		http.Error(w, errBody, http.StatusConflict)
		return q, false
	}

	return q, true
}

func displayMessages(w http.ResponseWriter, batch []queue.Message) {
	w.Header().Add("Content-Type", "text/csv")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("id,header,body,created_at,ttl\n"))
	for _, m := range batch {
		w.Write([]byte(fmt.Sprintf(
			"%d,%s,%s,%d,%d\n",
			m.ID, m.Header, m.Body, m.CreatedAt, m.TTL,
		)))
	}
}
