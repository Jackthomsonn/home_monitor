package functions

import (
	"encoding/json"
	"io"
	"net/http"

	"cloud.google.com/go/datastore"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"jackthomson.com/functions/models"
	"jackthomson.com/functions/utils"
)

type DiscoveryRequest struct {
	Devices   []models.Device `json:"devices"`
	Timestamp string          `json:"timestamp"`
	ClientID  string          `json:"client_id"`
}

type DiscoveryDataStore struct {
	Ip         []int  `json:"ip"`
	Alias      string `json:"alias"`
	Feature    string `json:"feature"`
	OnTime     int    `json:"on_time"`
	DeviceId   string `json:"device_id"`
	RelayState int    `json:"relay_state"`
	DeviceType string `json:"device_type"`
	ClientID   string `json:"client_id"`
}

func DiscoverDevices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	api_key := r.Header.Get("api_key")

	if err := utils.CheckApiKey(api_key); err != nil {
		utils.Logger().Error("Error checking API key", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
		w.Write([]byte("Error checking API key"))
		return
	}

	var data DiscoveryRequest

	body, err := io.ReadAll(r.Body)

	if err != nil {
		utils.Logger().Error("Error reading request body", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if err := json.Unmarshal(body, &data); err != nil {
		println(err.Error())
		utils.Logger().Error("Error unmarshalling request body", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		utils.Logger().Error("Error unmarshalling request body", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, device := range data.Devices {
		key := datastore.NameKey("Device", device.DeviceId, nil)
		_, err := utils.WriteToDatastore(key, &DiscoveryDataStore{
			Alias:      device.Alias,
			Ip:         device.Ip,
			Feature:    device.Feature,
			OnTime:     device.OnTime,
			DeviceId:   device.DeviceId,
			RelayState: device.RelayState,
			DeviceType: device.DeviceType,
			ClientID:   data.ClientID,
		})
		if err != nil {
			println(err.Error())
			utils.Logger().Error("Error writing to datastore", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
