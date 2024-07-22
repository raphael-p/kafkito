package resolvers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/raphael-p/kafkito/server/config"
	"github.com/raphael-p/kafkito/server/queue"
	"github.com/raphael-p/kafkito/server/utils"
)

var queues queue.QueueMap = make(queue.QueueMap)

func CreateQueue(w http.ResponseWriter, r *http.Request) {
	queueName, ok := parseQueueName(w, r, "name")
	if !ok {
		return
	}

	err := queues.AddQueue(queueName)
	if err != nil {
		errBody := "error adding queue: " + err.Error()
		http.Error(w, errBody, http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func DeleteQueue(w http.ResponseWriter, r *http.Request) {
	queueName, ok := parseQueueName(w, r, "name")
	if !ok {
		return
	}
	delete(queues, queueName)
}

func RenameQueue(w http.ResponseWriter, r *http.Request) {
	q, ok := getQueue(w, r, queues, "oldName")
	if !ok {
		return
	}
	newName, ok := parseQueueName(w, r, "newName")
	if !ok {
		return
	}

	q.Name = newName
	queues[q.Name] = q
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
	q, ok := getQueue(w, r, queues, "name")
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
			http.StatusConflict,
		)
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
	queues[q.Name] = q
	w.WriteHeader(http.StatusCreated)
}

func ReadMessages(w http.ResponseWriter, r *http.Request) {
	q, ok := getQueue(w, r, queues, "name")
	if !ok {
		return
	}

	cursorStr := r.URL.Query().Get("cursor")
	var cursorID uint64
	if cursorStr != "" {
		cursorInt, err := strconv.Atoi(cursorStr)
		if err != nil {
			errBody := "error parsing 'cursor' query param: " + err.Error()
			http.Error(w, errBody, http.StatusBadRequest)
			return
		}
		if cursorInt < 0 {
			errBody := "error parsing 'cursor' query param: value must be positive"
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

	displayMessages(w, batch)
}

func ConsumeMessage(w http.ResponseWriter, r *http.Request) {
	messageID, ok := parseMessageID(w, r)
	if !ok {
		return
	}

	for _, q := range queues {
		foundIndex := -1
		for index, m := range q.Messages {
			if m.ID == messageID {
				foundIndex = index
				break
			}
		}
		if foundIndex > -1 {
			displayMessages(w, []queue.Message{q.Messages[foundIndex]})
			q.Messages = utils.RemoveFromSlice(q.Messages, foundIndex)
			queues[q.Name] = q
			return
		}
	}

	errBody := fmt.Sprint("no message found with ID: ", messageID)
	http.Error(w, errBody, http.StatusNotFound)
}
