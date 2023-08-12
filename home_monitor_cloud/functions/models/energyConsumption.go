package models

import "time"

type EnergyConsumption struct {
	Timestamp time.Time `json:"timestamp"`
	PowerMw   float32   `json:"power_mw"`
	Alias     string    `json:"alias"`
}
