package queue

import (
	"errors"
	"fmt"
	"time"

	"github.com/raphael-p/kafkito/server/config"
	"github.com/raphael-p/kafkito/server/utils"
)

var messageAutoInc utils.AutoIncrement

type QueueMap map[string]Queue

type Queue struct {
	Name      string
	Messages  []Message
	CreatedAt int64
}
type Message struct {
	ID        uint64
	Header    string
	Body      string
	CreatedAt int64
	TTL       int64
}

func (queues QueueMap) AddQueue(newQueueName string) error {
	if len(queues) >= int(config.MAX_QUEUES) {
		return errors.New(fmt.Sprint(
			"too many queues, max is: ",
			config.MAX_QUEUES,
		))
	}

	if _, ok := queues[newQueueName]; ok {
		return errors.New(fmt.Sprint(
			"queue already exists with name: " + newQueueName,
		))
	}

	queues[newQueueName] = Queue{
		Name:      newQueueName,
		Messages:  make([]Message, 0, config.MAX_QUEUE_LENGTH),
		CreatedAt: time.Now().Unix(),
	}

	return nil
}

// remove messages from queue which have exceded their TTL
func (queues QueueMap) GetQueue(queueName string) (Queue, bool) {
	q, ok := queues[queueName]
	if !ok {
		return q, false
	}

	purgedMessages := make([]Message, 0, len(q.Messages))
	for _, message := range q.Messages {
		if message.CreatedAt+message.TTL > time.Now().Unix() {
			purgedMessages = append(purgedMessages, message)
		}
	}
	q.Messages = purgedMessages
	queues[q.Name] = q

	return q, true
}

func MakeMessage(header string, body string, ttl int64) (Message, error) {
	if header == "" || body == "" {
		return Message{}, errors.New("message header or body must not be empty")
	}

	if len(header) > int(config.MAX_MESSAGE_HEADER_BYTES) {
		return Message{}, errors.New(fmt.Sprint(
			"message header is too long, max is: ",
			config.MAX_MESSAGE_HEADER_BYTES,
		))
	}

	if len(body) > int(config.MAX_MESSAGE_BODY_BYTES) {
		return Message{}, errors.New(fmt.Sprint(
			"message body is too long, max is: ",
			config.MAX_MESSAGE_BODY_BYTES,
		))
	}

	var messageTTL int64
	if ttl == 0 {
		messageTTL = config.MESSAGE_TTL
	} else {
		messageTTL = ttl
	}

	return Message{
		ID:        messageAutoInc.ID(),
		Header:    header,
		Body:      body,
		CreatedAt: time.Now().Unix(),
		TTL:       messageTTL,
	}, nil
}
