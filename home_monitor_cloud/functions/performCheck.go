package functions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type PerformCheckRequest struct {
  IntervalInMinutes *int `json:"interval_in_minutes"`
}

type Carbonintensity struct {
  Index     string `json:"index"`
  Forecast  int    `json:"forecast"`
  Actual    int    `json:"actual"`
}

type CarbonintensityData struct {
  From      string         `json:"from"`
  To        string         `json:"to"`
  Intensity Carbonintensity `json:"intensity"`
}

type CarbonintensityResponse struct {
  Data []CarbonintensityData `json:"data"`
}

type Action string

const (
  TURN_ON       Action = "TURN_ON"
  TURN_OFF      Action = "TURN_OFF"
  MAYBE_TURN_ON Action = "MAYBE_TURN_ON"
)

type Response struct {
  Action  Action `json:"action"`
  Index   string `json:"index"`
  Forecast int    `json:"forecast"`
  Unit    string `json:"unit"`
  From    string `json:"from"`
  To      string `json:"to"`
}

func PerformCheck(w http.ResponseWriter, r *http.Request) {
  var req PerformCheckRequest
  decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)

  utc := time.Now().UTC().Truncate(time.Second)
  
  now := utc.Format(time.RFC3339)

  nowPlusIntervalInMinutes := utc.Add(time.Minute * time.Duration(*req.IntervalInMinutes)).Format(time.RFC3339)

  if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
  
  defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

  if req.IntervalInMinutes == nil {
    http.Error(w, errors.New("interval_in_minutes is required").Error(), http.StatusBadRequest)
    return
  }

  data, err := getCarbonIntensity(w, *req.IntervalInMinutes, now, nowPlusIntervalInMinutes)

  if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

  carbonintensityResponse, err := determineCarbonIntensity(data, now, nowPlusIntervalInMinutes)

  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(carbonintensityResponse)
}

func getTemperature() (int, error) {
  result, err := http.Get("https://api.open-meteo.com/v1/forecast?latitude=51.45&longitude=-3.18&hourly=apparent_temperature&current_weather=true&precipitation_unit=inch&timezone=Europe%2FLondon")

  body, err := ioutil.ReadAll(result.Body)

  var data map[string]interface{}

  err = json.Unmarshal(body, &data)

  defer result.Body.Close()

  if err != nil {
    return 0, err
  }

  currentWeather := data["current_weather"].(map[string]interface{})
  apparentTemperature := currentWeather["apparent_temperature"].(float64)

  return int(apparentTemperature), nil
}

func getCarbonIntensity(w http.ResponseWriter, intervalInMinutes int, now string, nowPlusIntervalInMinutes string) (CarbonintensityResponse, error) {  
  result, err := http.Get(fmt.Sprintf("https://api.carbonintensity.org.uk/intensity/%s/%s", now, nowPlusIntervalInMinutes))

  if err != nil {
		return CarbonintensityResponse{Data: nil}, err;
	}

  body, err := ioutil.ReadAll(result.Body)

  if err != nil {
		return CarbonintensityResponse{Data: nil}, err;
	}

  defer result.Body.Close()

  var data CarbonintensityResponse

  err = json.Unmarshal(body, &data)

  if err != nil {
    return CarbonintensityResponse{Data: nil}, err;
  }

  return data, nil
}

func determineCarbonIntensity(data CarbonintensityResponse, now string, nowPlusIntervalInMinutes string) (Response, error) {
  latestData := data.Data[0]

  intensity := latestData.Intensity

  efficientIntensities := []string{"very low", "low"}

  midEfficientIntensities := []string{"moderate"}

  response := Response{Index: intensity.Index, Forecast: intensity.Forecast, Unit: "gCO2/kWh", From: now, To: nowPlusIntervalInMinutes}

  if contains(efficientIntensities, intensity.Index) {
    response.Action = TURN_ON

    return response, nil
  }

  if contains(midEfficientIntensities, intensity.Index) {
    response.Action = MAYBE_TURN_ON

    return response, nil
  }

  response.Action = TURN_OFF
  
  return response, nil
}

func contains(slice []string, item string) bool {
  for _, s := range slice {
    if s == item {
      return true
    }
  }

  return false
}