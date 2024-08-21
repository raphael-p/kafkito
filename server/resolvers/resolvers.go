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
		writeError(w, errBody, http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	utils.LogInfo(fmt.Sprintf("queue created: %s\n", queueName))
}

func DeleteQueue(w http.ResponseWriter, r *http.Request) {
	queueName, ok := parseQueueName(w, r, "name")
	if !ok {
		return
	}

	delete(queues, queueName)
	utils.LogInfo(fmt.Sprintf("queue deleted: %s\n", queueName))
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

	oldName := q.Name
	q.Name = newName
	queues[q.Name] = q
	delete(queues, oldName)
	utils.LogInfo(fmt.Sprintf(
		"queue renamed from %s to %s\n",
		oldName, newName,
	))
}

func ListQueues(w http.ResponseWriter, r *http.Request) {
	if len(queues) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	writeQueuesCSV(w, queues)
}

func PublishMessage(w http.ResponseWriter, r *http.Request) {
	q, ok := getQueue(w, r, queues, "name")
	if !ok {
		return
	}

	if len(q.Messages) >= int(config.Values.MaxQueueLength) {
		errBody := fmt.Sprint(
			"too many messages in queue, max is: ",
			config.Values.MaxQueueLength,
		)
		writeError(w, errBody, http.StatusConflict)
		return
	}

	err := r.ParseForm()
	if err != nil {
		writeError(w, "error parsing request body", http.StatusBadRequest)
		return
	}
	header := r.FormValue("header")
	body := r.FormValue("body")

	message, err := queue.MakeMessage(header, body, config.Values.MessageTTL)
	if err != nil {
		errBody := "error creating message: " + err.Error()
		writeError(w, errBody, http.StatusBadRequest)
		return
	}

	q.Messages = append(q.Messages, message)
	queues[q.Name] = q
	w.WriteHeader(http.StatusCreated)
	utils.LogInfo(fmt.Sprintf(
		"published message %d with header %s to queue %s\n",
		message.ID, message.Header, q.Name,
	))
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
			writeError(w, errBody, http.StatusBadRequest)
			return
		}
		if cursorInt < 0 {
			errBody := "error parsing 'cursor' query param: value must be positive"
			writeError(w, errBody, http.StatusBadRequest)
			return
		}
		cursorID = uint64(cursorInt)
	}

	batch := make([]queue.Message, 0, config.Values.MessageBatchSize)
	for _, m := range q.Messages {
		if len(batch) >= int(config.Values.MessageBatchSize) {
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

	writeMessagesCSV(w, batch)
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
			writeMessagesCSV(w, []queue.Message{q.Messages[foundIndex]})
			q.Messages = utils.RemoveFromSlice(q.Messages, foundIndex)
			queues[q.Name] = q
			utils.LogInfo(fmt.Sprintf(
				"consumed message %d on queue %s\n",
				messageID, q.Name,
			))
			return
		}
	}

	errBody := fmt.Sprint("no message found with ID: ", messageID)
	writeError(w, errBody, http.StatusNotFound)
}
