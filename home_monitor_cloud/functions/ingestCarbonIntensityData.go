package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"jackthomson.com/functions/models"
	"jackthomson.com/functions/services"
	"jackthomson.com/functions/utils"
)

type CarbonIntensityIngestion struct {
	Actual    float64 `json:"actual"`
	ForeCast  float64 `json:"forecast"`
	Timestamp string  `json:"from"`
}

func IngestCarbonIntensityData(w http.ResponseWriter, r *http.Request) {
	utils.Logger().Info("IngestCarbonIntensityData", zap.Field{Key: "method", Type: zapcore.StringType, String: r.Method}, zap.Field{Key: "url", Type: zapcore.StringType, String: r.URL.String()})
	w.Header().Set("Content-Type", "application/json")

	utc := time.Now().UTC().Truncate(time.Second)

	now := utc.Format(time.RFC3339)
	previousDay := utc.AddDate(0, 0, -1).Format(time.RFC3339)

	data, err := getCarbonIntensity(w, previousDay, now)

	if err != nil {
		utils.Logger().Error("Error getting carbon intensity data", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Error{Message: err.Error()})
		return
	}

	ingestionValues := []CarbonIntensityIngestion{}

	for _, intensity := range data.Data {
		ingestionValues = append(ingestionValues, CarbonIntensityIngestion{
			Actual:    float64(intensity.Intensity.Actual),
			ForeCast:  float64(intensity.Intensity.Forecast),
			Timestamp: intensity.From,
		})
	}

	if err = services.InsertDataIntoBiqQuery(context.Background(), ingestionValues, "home_monitor_carbon_intensity"); err != nil {
		utils.Logger().Error("Error inserting carbon intensity data", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Error{Message: err.Error()})
		return
	}

	utils.Logger().Info("Successfully ingested carbon intensity data")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ingestionValues)
}

func getCarbonIntensity(_ http.ResponseWriter, previousDay string, now string) (models.CarbonintensityResponse, error) {
	utils.Logger().Info("GetCarbonIntensity", zap.Field{Key: "now", Type: zapcore.StringType, String: previousDay}, zap.Field{Key: "now", Type: zapcore.StringType, String: now})
	result, err := http.Get(fmt.Sprintf("https://api.carbonintensity.org.uk/intensity/%s/%s", previousDay, now))

	if err != nil {
		utils.Logger().Error("Error getting carbon intensity", zap.Error(err))
		return models.CarbonintensityResponse{Data: nil}, err
	}

	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)

	if err != nil {
		utils.Logger().Error("Error reading carbon intensity response", zap.Error(err))
		return models.CarbonintensityResponse{Data: nil}, err
	}

	var data models.CarbonintensityResponse

	err = json.Unmarshal(body, &data)

	if err != nil {
		utils.Logger().Error("Error unmarshalling carbon intensity response", zap.Error(err))
		return models.CarbonintensityResponse{Data: nil}, err
	}

	utils.Logger().Info("Successfully got carbon intensity")
	return data, nil
}
