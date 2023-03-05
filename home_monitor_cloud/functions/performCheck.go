package functions

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"jackthomson.com/functions/enums"
	"jackthomson.com/functions/models"
	"jackthomson.com/functions/services"
	"jackthomson.com/functions/utils"
)

type PerformCheckRequest struct {
	IntervalInMinutes *int `json:"interval_in_minutes"`
}

type Response struct {
	Intensities []models.Carbonintensity `json:"intensities"`
}

func PerformCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req PerformCheckRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Error{Message: err.Error()})
		return
	}

	defer r.Body.Close()

	utc := time.Now().UTC().Truncate(time.Second)

	now := utc.Format(time.RFC3339)

	nowPlusIntervalInMinutes := utc.Add(time.Minute * time.Duration(*req.IntervalInMinutes)).Format(time.RFC3339)

	if req.IntervalInMinutes == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Error{Message: errors.New("interval_in_minutes is required").Error()})
		return
	}

	data, err := services.GetCarbonIntensity(w, now, nowPlusIntervalInMinutes)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Error{Message: err.Error()})
		return
	}

	carbonintensityResponse, err := determineCarbonIntensity(data)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Error{Message: err.Error()})
		return
	}

	response := Response{Intensities: carbonintensityResponse}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func determineCarbonIntensity(data services.CarbonintensityResponse) ([]models.Carbonintensity, error) {
	latestData := data.Data

	efficientIntensities := []string{"very low", "low"}

	midEfficientIntensities := []string{"moderate"}

	response := []models.Carbonintensity{}

	for i := range latestData {
		if utils.Contains(efficientIntensities, latestData[i].Intensity.Index) {
			response = append(response, models.Carbonintensity{Index: latestData[i].Intensity.Index, Actual: latestData[i].Intensity.Actual, Forecast: latestData[i].Intensity.Forecast, Unit: "gCO2/kWh", From: latestData[i].From, To: latestData[i].To, Action: enums.TURN_ON})
			continue
		}

		if utils.Contains(midEfficientIntensities, latestData[i].Intensity.Index) {
			response = append(response, models.Carbonintensity{Index: latestData[i].Intensity.Index, Actual: latestData[i].Intensity.Actual, Forecast: latestData[i].Intensity.Forecast, Unit: "gCO2/kWh", From: latestData[i].From, To: latestData[i].To, Action: enums.MAYBE_TURN_ON})
			continue
		}

		response = append(response, models.Carbonintensity{Index: latestData[i].Intensity.Index, Actual: latestData[i].Intensity.Actual, Forecast: latestData[i].Intensity.Forecast, Unit: "gCO2/kWh", From: latestData[i].From, To: latestData[i].To, Action: enums.TURN_OFF})
	}

	return response, nil
}
