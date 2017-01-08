package requests

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/parnurzeal/gorequest"
)

var (
	apiKey = ""
)

type SensorReading struct {
	Device `json:"sensor"`
	Data   json.RawMessage `json:"data"`
}

type Device struct {
	Device string `json:"device"`
	Name   string `json:"name,omitempty"`
}

type APIKey struct {
	Token string `json:"token"`
}

type Payload struct {
	Token         string `json:"token"`
	SensorReading `json:"data"`
}

func SetKey(w http.ResponseWriter, r *http.Request) {
	log.Println("Setting API Key")
	defer r.Body.Close()
	if r.Method != "POST" {
		log.Println("Method not POST")
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var token APIKey
	err := json.NewDecoder(r.Body).Decode(&token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	apiKey = token.Token
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
	payload := Payload{apiKey, message}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	payload.Publish()
}

func (reading *SensorReading) Publish() error {
	data, err := json.Marshal(&reading)
	log.Println(string(data))
	request := gorequest.New()
	request.Post("https://smoker-relay.us/reading").Set("Notes", "gorequest").Send(string(data)).End(printStatus)

	/*req, err := http.NewRequest("POST", "https://smoker-relay.us/reading,", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	log.Println(err.Error())*/
	return err
}

func (payload *Payload) Publish() error {
	data, err := json.Marshal(&payload)
	log.Println(string(data))
	request := gorequest.New()
	request.Post("https://smoker-relay.us/reading").Set("Notes", "gorequest").Send(string(data)).End(printStatus)
	return err
}
func printStatus(resp gorequest.Response, body string, errs []error) {
	fmt.Println(resp.Status)
}
