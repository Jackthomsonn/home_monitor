package models

import "time"

type EnergyConsumption struct {
	Timestamp time.Time `json:"timestamp"`
	Power_Mw  float64   `json:"power_mw"`
	Alias     string    `json:"alias"`
}
