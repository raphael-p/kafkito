package resolvers

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/raphael-p/kafkito/cli/utils"
)

func CreateQueue(queueName string) {
	response := utils.KafkitoPost("/queue/" + queueName)
	if response.Error != nil {
		fmt.Println(response.Error.Error())
		return
	}

	fmt.Println("queue created:", queueName)
}

func RenameQueue(oldQueueName, newQueueName string) {
	response := utils.KafkitoPost(
		"/queue/" + oldQueueName + "/rename/" + newQueueName,
	)
	if response.Error != nil {
		fmt.Println(response.Error.Error())
		return
	}

	fmt.Printf(
		"queue renamed from %s to %s\n",
		oldQueueName, newQueueName,
	)
}

func DeleteQueue(queueName string) {
	response := utils.KakitoDelete("/queue/" + queueName)
	if response.Error != nil {
		fmt.Println(response.Error.Error())
		return
	}

	fmt.Println("queue deleted:", queueName)
}

func ListQueues() {
	response := utils.KafkitoGetCSV("/queues")
	if response.Error != nil {
		fmt.Println(response.Error.Error())
		return
	}

	if response.StatusCode == http.StatusNoContent {
		fmt.Println("there are no queues")
		return
	}

	nameWidth := int(math.Max(
		float64(utils.GetQueueNameMaxLength()),
		float64(len("name")),
	))
	countWidth := len("message_count")
	columnWidths := []int{nameWidth, countWidth, utils.TIME_CHARS}

	dataFormatter := func(index int, data string) (string, error) {
		switch index {
		case 2:
			unixSeconds, err := strconv.Atoi(data)
			if err != nil {
				return "", fmt.Errorf("error: %s", err)
			}
			datetime := time.Unix(int64(unixSeconds), 0)
			return utils.FormatTime(datetime), nil
		default:
			return data, nil
		}
	}

	displayCSV(response.BodyStream, columnWidths, dataFormatter)
}

func PublishMessage(queueName, header, body string) {
	response := utils.KafkitoPostForm("/queue/"+queueName+"/publish", url.Values{
		"header": {header},
		"body":   {body},
	})

	if response.Error != nil {
		fmt.Println(response.Error.Error())
		return
	}

	messageID, err := strconv.Atoi(response.BodyString)
	if err != nil {
		fmt.Println("error: message ID could not be parsed from response")
	}

	fmt.Printf(
		"published message %d with header %s to queue %s\n",
		messageID, header, queueName,
	)
}

func ReadMessages() {
	fmt.Print("placeholder for 'read' command\n")
}

func ConsumeMessage() {
	fmt.Print("placeholder for 'consume' command\n")
}
