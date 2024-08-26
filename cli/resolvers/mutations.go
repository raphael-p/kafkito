package resolvers

import (
	"fmt"

	"github.com/raphael-p/kafkito/cli/utils"
)

func handleFailure(response utils.KafkitoResponse) bool {
	if response.Error != nil {
		fmt.Println("error: kafkito is not running on port", utils.GetPort())
		return true
	}

	if !utils.IsSuccessful(response.StatusCode) {
		fmt.Printf(
			"error: status code %d: %s",
			response.StatusCode, response.Body,
		)
		return true
	}

	return false
}

func CreateQueue(queueName string) {
	response := utils.KafkitoPost("/queue/"+queueName, "", "")
	if handleFailure(response) {
		return
	}

	fmt.Println("queue created:", queueName)
}

func RenameQueue(oldQueueName, newQueueName string) {
	response := utils.KafkitoPost(
		"/queue/"+oldQueueName+"/rename/"+newQueueName,
		"",
		"",
	)
	if handleFailure(response) {
		return
	}

	fmt.Printf(
		"queue renamed from %s to %s\n",
		oldQueueName, newQueueName,
	)
}

func DeleteQueue(queueName string) {
	response := utils.KakitoDelete("/queue/" + queueName)
	if handleFailure(response) {
		return
	}

	fmt.Println("queue deleted:", queueName)
}

func ListQueues() {
	fmt.Print("placeholder for 'list' command\n")
}

func PublishQueue() {
	fmt.Print("placeholder for 'publish' command\n")
}

func ReadMessages() {
	fmt.Print("placeholder for 'read' command\n")
}

func ConsumeMessage() {
	fmt.Print("placeholder for 'consume' command\n")
}
