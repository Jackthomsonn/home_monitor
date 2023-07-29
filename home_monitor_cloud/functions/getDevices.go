package functions

import (
	"encoding/json"
	"log"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"jackthomson.com/functions/utils"
)

func GetDevices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := utils.CheckApiKey(r.Header.Get("api_key")); err != nil {
		utils.Logger().Error("Error checking API key", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	log.Println("Getting devices")
	keys, err := utils.ReadAllFromDataStore()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(keys)
}
