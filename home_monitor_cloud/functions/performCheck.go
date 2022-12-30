package functions

import (
	"encoding/json"
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

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

  if req.IntervalInMinutes == nil {
    http.Error(w, "interval_in_minutes is required", http.StatusBadRequest)
    return
  }

  utc := time.Now().UTC().Truncate(time.Second)
  
  now := utc.Format(time.RFC3339)

  nowPlusIntervalInMinutes := utc.Add(time.Minute * time.Duration(*req.IntervalInMinutes)).Format(time.RFC3339)

  result, err := http.Get(fmt.Sprintf("https://api.carbonintensity.org.uk/intensity/%s/%s", now, nowPlusIntervalInMinutes))

  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte(err.Error()))
    return
  }
  
  defer result.Body.Close()

  body, err := ioutil.ReadAll(result.Body)

  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte(err.Error()))
    return
  }

  var data CarbonintensityResponse

  err = json.Unmarshal(body, &data)

  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte(err.Error()))
    return
  }

  latestData := data.Data[0]

  intensity := latestData.Intensity

  efficientIntensities := []string{"very low", "low"}

  midEfficientIntensities := []string{"moderate"}

  if contains(efficientIntensities, intensity.Index) {
    response := Response{
      Action:  TURN_ON,
      Index:   intensity.Index,
      Forecast: intensity.Forecast,
      Unit:    "gCO2/kWh",
      From: now,
      To: nowPlusIntervalInMinutes,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
    return
  }

  if contains(midEfficientIntensities, intensity.Index) {
    response := Response{
      Action:  MAYBE_TURN_ON,
      Index:   intensity.Index,
      Forecast: intensity.Forecast,
      Unit:    "gCO2/kWh",
      From: now,
      To: nowPlusIntervalInMinutes,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
    return
  }

  response := Response{
    Action: TURN_OFF,
    Index:   intensity.Index,
    Forecast: intensity.Forecast,
    Unit:    "gCO2/kWh",
    From: now,
    To: nowPlusIntervalInMinutes,
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(response)
}

func contains(slice []string, item string) bool {
  for _, s := range slice {
    if s == item {
      return true
    }
  }

  return false
}