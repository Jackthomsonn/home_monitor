package functions

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"jackthomson.com/functions/config"
	"jackthomson.com/functions/models"
	"jackthomson.com/functions/services"
	"jackthomson.com/functions/strategies"
	"jackthomson.com/functions/utils"
)

func IngestConsumptionData(w http.ResponseWriter, r *http.Request) {
	utils.Logger().Info("IngestConsumptionData", zap.Field{Key: "method", Type: zapcore.StringType, String: r.Method}, zap.Field{Key: "url", Type: zapcore.StringType, String: r.URL.String()})
	w.Header().Set("Content-Type", "application/json")

	config := config.GetConfig()

	utils.Logger().Info("Using strategy", zap.Field{Key: "strategy", Type: zapcore.StringType, String: config.StrategyToUse})

	response, err := strategies.StrategyFactory(config.StrategyToUse)

	if err != nil {
		utils.Logger().Error("Error getting consumption data", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Error{Message: err.Error()})
		return
	}

	if response.Values == nil {
		utils.Logger().Info("No consumption data to ingest")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Values = response.Values[:len(response.Values)-1]

	bigqueryErr := services.InsertDataIntoBiqQuery(context.Background(), response.Values, "home_monitor_consumption_table")

	if bigqueryErr != nil {
		utils.Logger().Error("Error inserting consumption data", zap.Error(bigqueryErr))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Error{Message: err.Error()})
		return
	}

	utils.Logger().Info("Successfully ingested consumption data", zap.Field{Key: "start", Type: zapcore.StringType, String: response.Start}, zap.Field{Key: "end", Type: zapcore.StringType, String: response.End})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
