package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"jackthomson.com/functions/models"
)

type CarbonintensityResponse struct {
	Data []models.CarbonintensityData `json:"data"`
}

func GetCarbonIntensity(w http.ResponseWriter, now string, nowPlusIntervalInMinutes string) (CarbonintensityResponse, error) {
	result, err := http.Get(fmt.Sprintf("https://api.carbonintensity.org.uk/intensity/%s/%s", now, nowPlusIntervalInMinutes))

	if err != nil {
		return CarbonintensityResponse{Data: nil}, err
	}

	defer result.Body.Close()

	body, err := ioutil.ReadAll(result.Body)

	if err != nil {
		return CarbonintensityResponse{Data: nil}, err
	}

	var data CarbonintensityResponse

	err = json.Unmarshal(body, &data)

	if err != nil {
		return CarbonintensityResponse{Data: nil}, err
	}

	return data, nil
}
