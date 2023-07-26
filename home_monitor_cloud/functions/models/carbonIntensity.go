package models

type Carbonintensity struct {
	Index    string `json:"index,omitempty"`
	Forecast int    `json:"forecast,omitempty"`
	Actual   int    `json:"actual,omitempty"`
	Unit     string `json:"unit,omitempty"`
	From     string `json:"from,omitempty"`
	To       string `json:"to,omitempty"`
}
