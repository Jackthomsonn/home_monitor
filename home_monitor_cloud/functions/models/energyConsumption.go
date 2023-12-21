package models

import "time"

type EnergyConsumption struct {
	Timestamp    time.Time `json:"timestamp"`
	Power_Wh_Avg float64   `json:"power_wh_avg"`
	Alias        string    `json:"alias"`
}
