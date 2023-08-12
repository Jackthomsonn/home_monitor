package models

type ConsumptionResponse struct {
	Unit                string              `json:"unit"`
	Granularity         string              `json:"granularity"`
	Start               string              `json:"start"`
	End                 string              `json:"end"`
	ResponseTimestamp   string              `json:"responseTimestamp"`
	Resource            string              `json:"resource"`
	Values              []ConsumptionValues `json:"values"`
	AvailableCacheRange AvailableCacheRange `json:"availableCacheRange"`
}
