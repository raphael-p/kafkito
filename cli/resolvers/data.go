package resolvers

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

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
	response := utils.KafkitoGetStream("/queues")
	if response.Error != nil {
		fmt.Println(response.Error.Error())
		return
	}

	if response.StatusCode == http.StatusNoContent {
		fmt.Println("there are no queues")
		return
	}

	printRow := func(row string, isHeader bool) bool {
		cells := strings.Split(row, ",")
		if len(cells) != 3 {
			fmt.Println("error: expected 3 columns in CSV response, got", len(cells))
			return false
		}

		nameWidth := utils.CalculateWidth("name", int(utils.GetQueueNameMaxLength()))
		countWidth := utils.CalculateWidth("message_count", 0)
		dateTimeWidth := utils.CalculateWidth("created_at", utils.TIME_CHAR_COUNT)

		if isHeader {
			utils.PrintCell("Name", nameWidth)
			utils.PrintCell("Message Count", countWidth)
			utils.PrintCell("Created At", dateTimeWidth)
			return true
		}

		utils.PrintCell(cells[0], nameWidth)
		utils.PrintCell(cells[1], countWidth)

		createdTime, err := strconv.Atoi(cells[2])
		if err != nil {
			fmt.Println("error:", err.Error())
			return false
		}
		utils.PrintCell(utils.UnixToDateTime(createdTime), dateTimeWidth)
		return true
	}

	// TODO: !!!PAGING!!!
	displayCSV(response.BodyStream, printRow)
}

func ReadMessages(queueName string) {
	response := utils.KafkitoGetStream("/queue/" + queueName + "/messages")
	if response.Error != nil {
		fmt.Println(response.Error.Error())
		return
	}

	if response.StatusCode == http.StatusNoContent {
		fmt.Println("there are no messages for", queueName)
		return
	}

	printRow := func(row string, isHeader bool) bool {
		cells := strings.Split(row, ",")
		if len(cells) != 5 {
			fmt.Println("error: expected 5 columns in CSV response, got", len(cells))
			return false
		}

		idWidth := utils.CalculateWidth("id", -1)
		headerWidth := utils.CalculateWidth("header", int(utils.GetHeaderMaxLength()))
		bodyWidth := utils.CalculateWidth("body", utils.MAX_BODY_DISPLAY)
		dateTimeWidth := utils.CalculateWidth("created_at", utils.TIME_CHAR_COUNT)

		if isHeader {
			utils.PrintCell("ID", idWidth)
			utils.PrintCell("Header", headerWidth)
			utils.PrintCell("Body", bodyWidth)
			utils.PrintCell("Created At", dateTimeWidth)
			utils.PrintCell("Expires", dateTimeWidth)
			return true
		}

		utils.PrintCell(cells[0], idWidth)
		utils.PrintCell(cells[1], headerWidth)
		utils.PrintCell(cells[2], bodyWidth)

		createdTime, err := strconv.Atoi(cells[3])
		if err != nil {
			fmt.Println("error:", err.Error())
			return false
		}
		utils.PrintCell(utils.UnixToDateTime(createdTime), dateTimeWidth)

		expirySeconds, err := strconv.Atoi(cells[4])
		if err != nil {
			fmt.Println("error:", err.Error())
			return false
		}
		expiry := createdTime + expirySeconds
		utils.PrintCell(utils.UnixToDateTime(expiry), dateTimeWidth)

		return true
	}

	displayCSV(response.BodyStream, printRow)
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

func processMessageResponse(response utils.KafkitoResponse, messageID string) {
	if response.Error != nil {
		fmt.Println(response.Error.Error())
		return
	}

	if response.StatusCode == http.StatusNoContent {
		fmt.Println("no message found with ID: ", messageID)
		return
	}

	stream := response.BodyStream
	defer stream.Close()

	scanner := bufio.NewScanner(stream)
	headerRow := true
	for scanner.Scan() {
		if headerRow {
			headerRow = false
			continue
		}
		body := strings.Split(scanner.Text(), ",")[2]
		fmt.Println(body)
		return
	}

	fmt.Println("error: unexpected response") // unreachable
}

func ReadMessage(messageID string) {
	response := utils.KafkitoGetStream("/message/" + messageID)
	processMessageResponse(response, messageID)
}

func ConsumeMessage(messageID string) {
	response := utils.KakitoDeleteStream("/message/" + messageID)
	processMessageResponse(response, messageID)
}
