package functions

import (
	"encoding/json"
	"net/http"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"jackthomson.com/functions/models"
	"jackthomson.com/functions/utils"
)

type Command struct {
	DeviceIP   string `json:"device_ip"`
	Action     string `json:"action"`
	DeviceType string `json:"device_type"`
	DeviceID   string `json:"device_id"`
}

type CommandRequest struct {
	Action     string `json:"action"`
	DeviceIP   string `json:"device_ip"`
	DeviceType string `json:"device_type"`
	DeviceID   string `json:"device_id"`
}

func SendCommand(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	var originToUse string = "https://home-monitor.vercel.app"

	if os.Getenv("DEVELOPMENT_MODE") == "true" {
		originToUse = "http://localhost:5173"
	}

	w.Header().Set("Access-Control-Allow-Origin", originToUse)

	api_key := r.Header.Get("api_key")

	if err := utils.CheckApiKey(api_key); err != nil {
		utils.Logger().Error("Error checking API key", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
		w.Write([]byte("Error checking API key"))
		return
	}

	host, err := utils.GetSecret("projects/345305797254/secrets/emqx_host/versions/latest", r.Context())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Error{Message: err.Error()})
		return
	}

	opts := mqtt.NewClientOptions().AddBroker(host)

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
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Error{Message: token.Error().Error()})
		return
	}

	defer c.Disconnect(250)

	command := Command{DeviceIP: command_request.DeviceIP, Action: command_request.Action, DeviceType: command_request.DeviceType, DeviceID: command_request.DeviceID}
	jsonCommand, err := json.Marshal(command)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Error{Message: err.Error()})
		return
	}

	topic := "commands/" + command_request.DeviceID + "/ping"
	token := c.Publish(topic, 0, false, jsonCommand)
	token.Wait()
}
