package functions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

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
}

func PerformCheck(w http.ResponseWriter, r *http.Request) {
  utc := time.Now().UTC().Truncate(time.Second)
  
  now := utc.Format(time.RFC3339)

  nowPlus30Minutes := utc.Add(time.Minute * 30).Format(time.RFC3339)

  result, err := http.Get(fmt.Sprintf("https://api.carbonintensity.org.uk/intensity/%s/%s", now, nowPlus30Minutes))
  
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