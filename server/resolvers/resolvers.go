package resolvers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/raphael-p/kafkito/server/config"
	"github.com/raphael-p/kafkito/server/queue"
)

var queues queue.QueueMap = make(queue.QueueMap)

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

func getQueue(w http.ResponseWriter, r *http.Request, queues queue.QueueMap) (queue.Queue, bool) {
	var q queue.Queue

	queueName, ok := parseQueueName(w, r)
	if !ok {
		return q, false
	}

	q, ok = queues.GetQueue(queueName)
	if !ok {
		errBody := "error fetching queue: no queue with name: " + queueName
		http.Error(w, errBody, http.StatusBadRequest)
		return q, false
	}

	return q, true
}

func CreateQueue(w http.ResponseWriter, r *http.Request) {
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

func ListQueues(w http.ResponseWriter, r *http.Request) {
	if len(queues) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Add("Content-Type", "text/csv")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("name,created_at\n"))
	for _, q := range queues {
		w.Write([]byte(fmt.Sprint(q.Name, ",", q.CreatedAt, "\n")))
	}
}

func PublishMessage(w http.ResponseWriter, r *http.Request) {
	q, ok := getQueue(w, r, queues)
	if !ok {
		return
	}

	if len(q.Messages) >= int(config.MAX_QUEUE_LENGTH) {
		http.Error(
			w,
			fmt.Sprint(
				"too many messages in queue, max is: ",
				config.MAX_QUEUE_LENGTH,
			),
			http.StatusBadRequest,
		)
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
	queues[q.Name] = q
	w.WriteHeader(http.StatusCreated)
}

func ReadMessages(w http.ResponseWriter, r *http.Request) {
	q, ok := getQueue(w, r, queues)
	if !ok {
		return
	}

	cursorStr := r.URL.Query().Get("cursor")
	var cursorID uint64
	if cursorStr != "" {
		cursorInt, err := strconv.Atoi(r.URL.Query().Get("cursor"))
		if err != nil {
			errBody := "error parsing 'cursor' query param: " + err.Error()
			http.Error(w, errBody, http.StatusBadRequest)
			return
		}
		if cursorInt <= 0 {
			errBody := "error parsing 'cursor' query param: value must be greater than zero"
			http.Error(w, errBody, http.StatusBadRequest)
			return
		}
		cursorID = uint64(cursorInt)
	}

	batch := make([]queue.Message, 0, config.MESSAGE_BATCH_SIZE)
	for _, m := range q.Messages {
		if len(batch) >= int(config.MESSAGE_BATCH_SIZE) {
			break
		}
		if m.ID > cursorID {
			batch = append(batch, m)
		}
	}
	newCursor := cursorID
	if len(batch) > 0 {
		newCursor = batch[len(batch)-1].ID
	}
	w.Header().Add("X-Cursor", fmt.Sprintf("%d", newCursor))

	if len(batch) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

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
