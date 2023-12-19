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

func AggregateHomeTotals(w http.ResponseWriter, r *http.Request) {
	utils.Logger().Info("AggregateHomeTotals", zap.Field{Key: "method", Type: zapcore.StringType, String: r.Method}, zap.Field{Key: "url", Type: zapcore.StringType, String: r.URL.String()})
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()

	key, err := aggregateHomeTotals()
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

func aggregateHomeTotals() (datastore.Key, error) {
	utils.Logger().Info("Ingesting home totals", zap.Field{Key: "function", Type: zapcore.ReflectType, Interface: "AggregateHomeTotals"})

	client, err := bigquery.NewClient(context.Background(), "home-monitor-373013")

	if err != nil {
		utils.Logger().Error("Error creating bigquery client", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
		return datastore.Key{}, err
	}

	defer client.Close()

	queryStr := `
WITH
Last24Hours AS (
SELECT
    timestamp AS table1_timestamp,
    CAST(REGEXP_EXTRACT(payload, r'current_ma: (\d+)') AS INT64) AS current_ma,
    CAST(REGEXP_EXTRACT(payload, r'voltage_mv: (\d+)') AS INT64) AS voltage_mV,
    (CAST(REGEXP_EXTRACT(payload, r'current_ma: (\d+)') AS INT64) * CAST(REGEXP_EXTRACT(payload, r'voltage_mv: (\d+)') AS INT64)) / 1000000 AS power_mw
FROM
    ` + "`home-monitor-373013.home_monitor_dataset.home_monitor_table`" + `
WHERE
    TIMESTAMP_DIFF(CURRENT_TIMESTAMP(), timestamp, HOUR) <= 24
    AND REGEXP_EXTRACT(payload, r'current_ma: (\d+)') IS NOT NULL
    AND REGEXP_EXTRACT(payload, r'voltage_mv: (\d+)') IS NOT NULL ),
HourlyConsumption AS (
SELECT
    TIMESTAMP_TRUNC(table1_timestamp, HOUR) AS hour_start,
    SUM(power_mW) AS consumption_mw
FROM
    Last24Hours
GROUP BY
    hour_start ),
CarbonCalculation AS (
SELECT
    h.hour_start,
    h.consumption_mw,
    c.actual AS carbonIntensity
FROM
    HourlyConsumption h
JOIN
    ` + "`home-monitor-373013.home_monitor_dataset.home_monitor_carbon_intensity`" + ` c
ON
    h.hour_start = c.timestamp )
SELECT
    ROUND(SUM(consumption_mw) / 1000000, 2) AS consumption,
    -- Convert from mWh to kWh
    ROUND(SUM((consumption_mw / 1000000) * carbonIntensity), 2) AS carbonIntensity
FROM
    CarbonCalculation;
`

	query := client.Query(queryStr)

	it, err := query.Read(context.Background())

	if err != nil {
		utils.Logger().Error("Error reading bigquery data", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
		return datastore.Key{}, err
	}

	var row RowResponse

	err = it.Next(&row)

	if err == iterator.Done {
		utils.Logger().Error("Error reading bigquery data", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
		return datastore.Key{}, err
	}

	key, err := services.WriteToDatastore(datastore.NameKey("Total", "total", nil), &UserTotalsResponse{CarbonTotal: row.CarbonIntensity, ConsumptionTotal: row.Consumption})

	if err != nil {
		utils.Logger().Error("Error writing to datastore", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
		return datastore.Key{}, err
	}

	utils.Logger().Info("Successfully ingested home totals", zap.Field{Key: "function", Type: zapcore.ReflectType, Interface: "AggregateHomeTotals"})

	if err := services.RemoveDataFromRedis("Total"); err != nil {
		utils.Logger().Error("Error removing data from redis", zap.Error(err))
	}

	return key, nil
}
