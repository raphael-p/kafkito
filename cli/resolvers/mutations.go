package resolvers

import (
	"fmt"
	"net/http"

	"github.com/raphael-p/kafkito/cli/utils"
)

func handleFailure(response utils.KafkitoResponse, successStatus int) bool {
	if response.Error != nil {
		fmt.Println("error: kafkito is not running on port", utils.GetPort())
		return true
	}

	if response.StatusCode != successStatus {
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
	if handleFailure(response, http.StatusCreated) {
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
	if handleFailure(response, http.StatusOK) {
		return
	}

	fmt.Printf(
		"queue renamed from %s to %s\n",
		oldQueueName, newQueueName,
	)
}

func DeleteQueue() {
	fmt.Print("placeholder for 'delete' command\n")
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
