package models

type Device struct {
	Ip       []int  `json:"ip"`
	Alias    string `json:"alias"`
	Feature  string `json:"feature"`
	OnTime   int    `json:"on_time"`
	DeviceId string `json:"device_id"`
}
