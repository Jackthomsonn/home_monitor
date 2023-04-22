package functions

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"jackthomson.com/functions/models"
	"jackthomson.com/functions/utils"
)

func TriggerConsumptionData(w http.ResponseWriter, r *http.Request) {
	utils.Logger().Info("TriggerConsumptionData", zap.Field{Key: "method", Type: zapcore.StringType, String: r.Method}, zap.Field{Key: "url", Type: zapcore.StringType, String: r.URL.String()})
	w.Header().Set("Content-Type", "application/json")

	var data Data
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		utils.Logger().Error("Error decoding request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Error{Message: err.Error()})
		return
	}

	defer r.Body.Close()

	id, err := utils.PublishDataToGCP(nil, "consumption-ingestion")
	if err != nil {
		utils.Logger().Error("Error publishing data to GCP", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Error{Message: err.Error()})
		return
	}

	utils.Logger().Info("Successfully published data to GCP")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseType{Type: "Success", Id: id})
}
