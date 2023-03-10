package functions

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

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
	w.Header().Set("Content-Type", "application/json")

	utc := time.Now().UTC().Truncate(time.Second)

	now := utc.Format(time.RFC3339)
	previousDay := utc.AddDate(0, 0, -1).Format(time.RFC3339)

	data, err := services.GetCarbonIntensity(w, previousDay, now)

	if err != nil {
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

	if err = utils.InsertDataIntoBiqQuery(context.Background(), ingestionValues, "home_monitor_carbon_intensity"); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Error{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ingestionValues)
}
