package models

type HomeTotal struct {
	ConsumptionTotal float64 `json:"consumptionTotal,omitempty"`
	CarbonTotal      float64 `json:"carbonTotal,omitempty"`
}
