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
}

func DiscoverDevices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
		_, err := utils.WriteToDatastore(key, &device)
		if err != nil {
			utils.Logger().Error("Error writing to datastore", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
