package functions

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

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

	var res []models.EnergyConsumption = make([]models.EnergyConsumption, 0)

	if err := services.GetDataFromRedis(&res, "energy_consumption"); err == nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
		return
	}

	query := `SELECT REGEXP_EXTRACT(payload, r'alias: ([^,}]+)') AS alias, CEIL(AVG(CAST(REGEXP_EXTRACT(payload, r'power_mw: (\d+)') as INT64) )) as power_wh_avg FROM ` + "`home-monitor-373013.home_monitor_dataset.home_monitor_table`" + ` WHERE timestamp >= TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 1 HOUR) GROUP BY alias`

	response, err := services.GetDataFromBigQuery[models.EnergyConsumption](context.Background(), "home_monitor_table", query)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	delta := time.Duration(1) * time.Hour
	if err := services.CreateDataInRedis(response, "energy_consumption", delta); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
