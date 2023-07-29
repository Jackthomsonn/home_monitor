package functions

import (
	"encoding/json"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"jackthomson.com/functions/models"
)

type Command struct {
	DeviceIP   string `json:"device_ip"`
	Action     string `json:"action"`
	DeviceType string `json:"device_type"`
}

type CommandRequest struct {
	Action     string `json:"action"`
	DeviceIP   string `json:"device_ip"`
	DeviceType string `json:"device_type"`
}

func SendCommand(w http.ResponseWriter, r *http.Request) {
	opts := mqtt.NewClientOptions().AddBroker("35.187.59.21:1883")

	var command_request CommandRequest
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&command_request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Error{Message: err.Error()})
		return
	}

	defer r.Body.Close()

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	defer c.Disconnect(250)

	command := Command{DeviceIP: command_request.DeviceIP, Action: command_request.Action, DeviceType: command_request.DeviceType}
	jsonCommand, err := json.Marshal(command)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Error{Message: err.Error()})
		return
	}
	token := c.Publish("commands/host/test", 0, false, jsonCommand)
	token.Wait()
}
