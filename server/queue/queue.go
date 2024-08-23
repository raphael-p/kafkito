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

var QueueCSVHeader = []byte("name,message_count,created_at\n")

func (q Queue) ToCSVRow() string {
	return fmt.Sprintf(
		"%s,%d,%d\n",
		q.Name, len(q.Messages), q.CreatedAt,
	)
}

type Message struct {
	ID        uint64
	Header    string
	Body      string
	CreatedAt int64
	TTL       int64
}

var MessageCSVHeader = []byte("id,header,body,created_at,ttl\n")

func (m Message) ToCSVRow() string {
	return fmt.Sprintf(
		"%d,%s,%s,%d,%d\n",
		m.ID, m.Header, m.Body, m.CreatedAt, m.TTL,
	)
}

func (queues QueueMap) AddQueue(newQueueName string) error {
	if _, ok := queues[newQueueName]; ok {
		return errors.New(fmt.Sprint(
			"queue already exists with name: " + newQueueName,
		))
	}

	if len(queues) >= int(config.Values.MaxQueues) {
		return errors.New(fmt.Sprint(
			"too many queues, max is: ",
			config.Values.MaxQueues,
		))
	}

	queues[newQueueName] = Queue{
		Name:      newQueueName,
		Messages:  make([]Message, 0, config.Values.MaxQueueLength),
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

	if len(header) > int(config.Values.MaxMessageHeaderBytes) {
		return Message{}, errors.New(fmt.Sprint(
			"message header is too long, max is: ",
			config.Values.MaxMessageHeaderBytes,
		))
	}

	if err := utils.CheckNameFormat(header, "message header "); err != nil {
		return Message{}, err
	}

	if len(body) > int(config.Values.MaxMessageBodyBytes) {
		return Message{}, errors.New(fmt.Sprint(
			"message body is too long, max is: ",
			config.Values.MaxMessageBodyBytes,
		))
	}

	var messageTTL int64
	if ttl == 0 {
		messageTTL = config.Values.MessageTTL
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
