package functions

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"jackthomson.com/functions/models"
	"jackthomson.com/functions/services"
	"jackthomson.com/functions/utils"
)

func GetEnergyConsumption(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", "*")

	api_key := r.Header.Get("api_key")

	if err := utils.CheckApiKey(api_key); err != nil {
		utils.Logger().Error("Error checking API key", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(err)
		return
	}

	query := `SELECT REGEXP_EXTRACT(payload, r'alias: ([^,}]+)') AS alias, REGEXP_EXTRACT(payload, r'power_mw: (\d+)') AS power_mw, timestamp FROM ` + "`home-monitor-373013.home_monitor_dataset.home_monitor_table`" + ` WHERE timestamp >= TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 1 HOUR) ORDER BY timestamp DESC`

	response, err := services.GetDataFromBigQuery[models.EnergyConsumption](context.Background(), "home_monitor_table", query)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
