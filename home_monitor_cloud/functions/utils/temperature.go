package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type WeatherData struct {
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	GenerationTime float64 `json:"generationtime_ms"`
	UtcOffset      int     `json:"utc_offset_seconds"`
	Timezone       string  `json:"timezone"`
	TimezoneAbbr   string  `json:"timezone_abbreviation"`
	Elevation      float64 `json:"elevation"`
	CurrentWeather struct {
		Temperature  float64 `json:"temperature"`
		Windspeed    float64 `json:"windspeed"`
		Winddirection float64 `json:"winddirection"`
		Weathercode  int     `json:"weathercode"`
		Time         string  `json:"time"`
	} `json:"current_weather"`
	HourlyUnits struct {
		Time                string `json:"time"`
		ApparentTemperature string `json:"apparent_temperature"`
	} `json:"hourly_units"`
	Hourly struct {
		Time []string `json:"time"`
	} `json:"hourly"`
}

func GetTemperature() (float64, error) {
  result, err := http.Get("https://api.open-meteo.com/v1/forecast?latitude=51.45&longitude=-3.18&hourly=apparent_temperature&current_weather=true&precipitation_unit=inch&timezone=Europe%2FLondon")

  if err != nil {
    return 0, err
  }
  
  defer result.Body.Close()

  body, err := ioutil.ReadAll(result.Body)

  if err != nil {
    return 0, err
  }

  var data = WeatherData{}

  err = json.Unmarshal(body, &data)

  if err != nil {
    return 0, err
  }

  current_weather := data.CurrentWeather

  temp := current_weather.Temperature

  return temp, nil
}