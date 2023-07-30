package functions

import (
	"encoding/json"
	"log"
	"net/http"

	"cloud.google.com/go/datastore"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"jackthomson.com/functions/utils"
)

type Device struct {
	A          int
	K          *datastore.Key `datastore:"__key__"`
	Ip         []int          `json:"ip"`
	Alias      string         `json:"alias"`
	Feature    string         `json:"feature"`
	OnTime     int            `json:"on_time"`
	RelayState int            `json:"relay_state"`
	DeviceId   string         `json:"device_id"`
	DeviceType string         `json:"device_type"`
	ClientID   string         `json:"client_id"`
}

func GetDevices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	api_key := r.Header.Get("api_key")

	if err := utils.CheckApiKey(api_key); err != nil {
		utils.Logger().Error("Error checking API key", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
		w.Write([]byte("Error checking API key"))
		return
	}

	log.Println("Getting devices")

	var devices []Device

	if err := utils.ReadAllFromDataStore("Device", &devices); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(devices)
}
