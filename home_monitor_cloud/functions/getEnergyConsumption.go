package functions

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"jackthomson.com/functions/models"
	"jackthomson.com/functions/services"
)

func GetEnergyConsumption(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	var originToUse string = "https://home-monitor.vercel.app"

	if os.Getenv("DEVELOPMENT_MODE") == "true" {
		originToUse = "http://localhost:5173"
	}

	w.Header().Set("Access-Control-Allow-Origin", originToUse)

	query := `SELECT REGEXP_EXTRACT(payload, r'alias: ([^,}]+)') AS alias, SAFE_CAST(REGEXP_EXTRACT(payload, r'power_mw: (\d+)') AS FLOAT64) AS power_mw, timestamp FROM ` + "`home-monitor-373013.home_monitor_dataset.home_monitor_table`" + ` WHERE timestamp >= TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 1 HOUR) ORDER BY timestamp DESC`

	response, err := services.GetDataFromBigQuery[models.EnergyConsumption](context.Background(), "home_monitor_table", query)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
