package resolvers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

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

	columnWidths := []int{
		utils.CalculateWidth("name", int(utils.GetQueueNameMaxLength())),
		utils.CalculateWidth("message_count", 0),
		utils.CalculateWidth("created_at", utils.TIME_CHAR_COUNT),
	}

	dataFormatter := func(index int, data string) (string, error) {
		switch index {
		case 2:
			return utils.UnixToDateTime(data)
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

func ReadMessages(queueName string) {
	response := utils.KafkitoGetCSV("/queue/" + queueName + "/messages")
	if response.Error != nil {
		fmt.Println(response.Error.Error())
		return
	}

	if response.StatusCode == http.StatusNoContent {
		fmt.Println("there are no messages for", queueName)
		return
	}

	columnWidths := []int{
		utils.CalculateWidth("id", -1),
		utils.CalculateWidth("header", int(utils.GetHeaderMaxLength())),
		utils.CalculateWidth("body", utils.MAX_BODY_DISPLAY),
		utils.CalculateWidth("created_at", utils.TIME_CHAR_COUNT),
		utils.CalculateWidth("ttl", -1),
	}

	dataFormatter := func(index int, data string) (string, error) {
		switch index {
		case 2:
			return utils.TruncateString(data, utils.MAX_BODY_DISPLAY), nil
		case 3:
			return utils.UnixToDateTime(data)
		default:
			return data, nil
		}
	}

	displayCSV(response.BodyStream, columnWidths, dataFormatter)
}

func ConsumeMessage() {
	fmt.Print("placeholder for 'consume' command\n")
}
