package functions

import (
	"context"
	"encoding/json"
	"net/http"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/datastore"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/api/iterator"
	"jackthomson.com/functions/models"
	"jackthomson.com/functions/services"
	"jackthomson.com/functions/utils"
)

type UserTotalsResponse struct {
	ConsumptionTotal float64 `json:"consumptionTotal,omitempty"`
	CarbonTotal      float64 `json:"carbonTotal,omitempty"`
}

type RowResponse struct {
	CarbonIntensity float64 `json:"carbonIntensity,omitempty"`
	Consumption     float64 `json:"consumption,omitempty"`
}

func IngestHomeTotals(w http.ResponseWriter, r *http.Request) {
	utils.Logger().Info("IngestHomeTotals", zap.Field{Key: "method", Type: zapcore.StringType, String: r.Method}, zap.Field{Key: "url", Type: zapcore.StringType, String: r.URL.String()})
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()

	key, err := ingestHomeTotals()
	if err != nil {
		utils.Logger().Error("Error ingesting home totals", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Error{Message: err.Error()})
		return
	}

	utils.Logger().Info("Successfully ingested home totals")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseType{Type: "Success", Id: key.String()})
}

func ingestHomeTotals() (datastore.Key, error) {
	utils.Logger().Info("Ingesting home totals", zap.Field{Key: "function", Type: zapcore.ReflectType, Interface: "IngestHomeTotals"})

	client, err := bigquery.NewClient(context.Background(), "home-monitor-373013")

	if err != nil {
		utils.Logger().Error("Error creating bigquery client", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
		return datastore.Key{}, err
	}

	defer client.Close()

	query := client.Query("SELECT `home-monitor-373013.home_monitor_dataset.home_monitor_carbon_intensity`.timestamp as ts, MAX(`home-monitor-373013.home_monitor_dataset.home_monitor_carbon_intensity`.actual) as carbonIntensity, MAX(`home_monitor_dataset.home_monitor_consumption_table`.value) as consumption, FROM `home-monitor-373013.home_monitor_dataset.home_monitor_carbon_intensity` INNER JOIN `home_monitor_dataset.home_monitor_consumption_table` ON `home-monitor-373013.home_monitor_dataset.home_monitor_carbon_intensity`.timestamp = `home_monitor_dataset.home_monitor_consumption_table`.timestamp WHERE DATE(`home-monitor-373013.home_monitor_dataset.home_monitor_carbon_intensity`.timestamp) = DATE_SUB(CURRENT_DATE(), INTERVAL 1 DAY) AND Date(`home-monitor-373013.home_monitor_dataset.home_monitor_carbon_intensity`.timestamp) >= DATE_TRUNC(DATE_SUB(CURRENT_DATE(), INTERVAL 1 DAY), DAY) GROUP BY ts ORDER BY ts ASC")

	it, err := query.Read(context.Background())

	if err != nil {
		utils.Logger().Error("Error reading bigquery data", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
		return datastore.Key{}, err
	}

	if it.TotalRows == 0 {
		utils.Logger().Info("No data to ingest")
		return datastore.Key{}, nil
	}

	var row RowResponse
	err = it.Next(&row)

	if err == iterator.Done {
		utils.Logger().Error("Error reading bigquery data", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
		return datastore.Key{}, err
	}

	var carbonTotal float64
	var consumptionTotal float64

	for {
		var row RowResponse
		if err := it.Next(&row); err == iterator.Done {
			break
		}

		if err != nil {
			return datastore.Key{}, err
		}

		carbonTotal += row.CarbonIntensity * row.Consumption
		consumptionTotal += row.Consumption
	}

	carbonTotal = float64(int(carbonTotal*100)) / 100
	consumptionTotal = float64(int(consumptionTotal*100)) / 100

	key, err := services.WriteToDatastore(datastore.NameKey("Total", "total", nil), &UserTotalsResponse{CarbonTotal: carbonTotal, ConsumptionTotal: consumptionTotal})

	if err != nil {
		utils.Logger().Error("Error writing to datastore", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
		return datastore.Key{}, err
	}

	utils.Logger().Info("Successfully ingested home totals", zap.Field{Key: "function", Type: zapcore.ReflectType, Interface: "IngestHomeTotals"})

	if err := services.RemoveDataFromRedis("Total"); err != nil {
		utils.Logger().Error("Error removing data from redis", zap.Error(err))
	}

	return key, nil
}
