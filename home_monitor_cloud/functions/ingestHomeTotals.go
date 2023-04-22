package functions

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"jackthomson.com/functions/models"
	"jackthomson.com/functions/services"
	"jackthomson.com/functions/utils"
)

func IngestHomeTotals(w http.ResponseWriter, r *http.Request) {
	utils.Logger().Info("IngestHomeTotals", zap.Field{Key: "method", Type: zapcore.StringType, String: r.Method}, zap.Field{Key: "url", Type: zapcore.StringType, String: r.URL.String()})
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()

	key, err := services.IngestHomeTotals()
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
