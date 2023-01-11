package models

type CarbonintensityData struct {
  From      string         `json:"from"`
  To        string         `json:"to"`
  Intensity Carbonintensity `json:"intensity"`
}