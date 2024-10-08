package resolvers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/raphael-p/kafkito/server/config"
	"github.com/raphael-p/kafkito/server/queue"
	"github.com/raphael-p/kafkito/server/utils"
)

func parseQueueName(w http.ResponseWriter, r *http.Request, pattern string) (string, bool) {
	queueName := strings.TrimSpace(r.PathValue(pattern))
	errPrefix := "error parsing queue name: "

	if queueName == "" {
		errBody := errPrefix + "queue name must not be empty"
		writeError(w, errBody, http.StatusBadRequest)
		return "", false
	}

	if err := utils.CheckNameFormat(queueName, errPrefix); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return "", false
	}

	if len(queueName) > int(config.Values.MaxQueueNameBytes) {
		errBody := fmt.Sprint(
			errPrefix,
			"queue name is too long, max is: ",
			config.Values.MaxQueueNameBytes,
		)
		writeError(w, errBody, http.StatusBadRequest)
		return "", false
	}

	return queueName, true
}

func parseMessageID(w http.ResponseWriter, r *http.Request) (uint64, bool) {
	messageIDStr := r.PathValue("id")
	errPrefix := "error parsing message ID: "

	if messageIDStr == "" {
		errBody := errPrefix + "must not be empty"
		writeError(w, errBody, http.StatusBadRequest)
		return 0, false
	}

	messageIDInt, err := strconv.Atoi(messageIDStr)
	if err != nil {
		errBody := errPrefix + err.Error()
		writeError(w, errBody, http.StatusBadRequest)
		return 0, false
	}
	if messageIDInt <= 0 {
		errBody := errPrefix + "must be greater than zero"
		writeError(w, errBody, http.StatusBadRequest)
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
		writeError(w, errBody, http.StatusConflict)
		return q, false
	}

	return q, true
}

func writeMessagesCSV(w http.ResponseWriter, messages []queue.Message) {
	w.Header().Add("content-type", "text/csv")
	w.WriteHeader(http.StatusOK)
	w.Write(queue.MessageCSVHeader)
	for _, m := range messages {
		w.Write([]byte(m.ToCSVRow()))
	}
}

func writeQueuesCSV(w http.ResponseWriter, queues queue.QueueMap) {
	w.Header().Add("content-type", "text/csv")
	w.WriteHeader(http.StatusOK)
	w.Write(queue.QueueCSVHeader)
	for _, q := range queues {
		w.Write([]byte(q.ToCSVRow()))
	}
}

func writeError(w http.ResponseWriter, errBody string, code int) {
	utils.LogWarning(errBody)
	http.Error(w, errBody, code)
}
