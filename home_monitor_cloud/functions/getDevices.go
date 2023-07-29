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
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	api_key := r.Header.Get("api_key")

	if err := utils.CheckApiKey(api_key); err != nil {
		utils.Logger().Error("Error checking API key", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
		w.Write([]byte("Error checking API key"))
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
