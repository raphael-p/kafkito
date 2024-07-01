package resolvers

import (
	"fmt"
	"net/http"

	"github.com/raphael-p/kafkito/server/queue"
)

func CreateQueue(queues queue.QueueMap) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// body, err := io.ReadAll(r.Body)
		// if err != nil {
		// 	http.Error(w, "error reading request body", http.StatusInternalServerError)
		// 	return
		// }

		queueName := r.PathValue("name")
		// TOOD: error handling -> empty queue and long queue

		queue.AddQueue(queues, queueName) // TODO: error handling

		fmt.Println(queues)
		w.WriteHeader(http.StatusCreated)
	}
}
