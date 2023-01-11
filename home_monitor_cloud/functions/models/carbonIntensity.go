package models

import (
	"jackthomson.com/functions/enums"
)

type Carbonintensity struct {
  Action  enums.Action `json:"action,omitempty"`
  Index   string `json:"index,omitempty"`
  Forecast int    `json:"forecast,omitempty"`
  Actual   int    `json:"actual,omitempty"`
  Unit    string `json:"unit,omitempty"`
  From    string `json:"from,omitempty"`
  To      string `json:"to,omitempty"`
}