package queue

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/raphael-p/kafkito/server/config"
)

type QueueMap map[string]Queue

type Queue struct {
	Name      string
	Messages  []Message
	CreatedAt int64
}

type MessageBody [config.MAX_MESSAGE_BODY_BYTES]byte

type Message struct {
	UUID      uint32
	Header    string
	Body      MessageBody
	CreatedAt int64
	TTL       uint32
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

func MakeMessage(header string, body string, ttl uint32) (Message, error) {
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
	var messageBody MessageBody
	for i := 0; i < len(body); i++ {
		messageBody[i] = body[i]
	}

	var messageTTL uint32
	if ttl == 0 {
		messageTTL = config.MESSAGE_TTL
	} else {
		messageTTL = ttl
	}

	return Message{
		UUID:      uuid.New().ID(),
		Header:    header,
		Body:      messageBody,
		CreatedAt: time.Now().Unix(),
		TTL:       messageTTL,
	}, nil
}
