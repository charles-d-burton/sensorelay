package requests

import (
	"encoding/json"
	"log"
	"net/http"
)

type SensorReading struct {
	Sensor string          `json:"sensor"`
	Topic  string          `json:"topic"`
	Data   json.RawMessage `json:"data"`
}

//Reading ... process sensor data input
func Reading(w http.ResponseWriter, r *http.Request) {
	log.Println("Request Received: ")

	defer r.Body.Close()
	if r.Method != "POST" {
		log.Println("Method not POST")
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var message SensorReading

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	data, err := json.Marshal(&message.Data)

	log.Println(message.Sensor)
	log.Println(message.Topic)
	log.Println(string(data))
}
