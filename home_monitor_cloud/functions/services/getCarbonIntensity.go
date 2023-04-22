package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"jackthomson.com/functions/models"
	"jackthomson.com/functions/utils"
)

type CarbonintensityResponse struct {
	Data []models.CarbonintensityData `json:"data"`
}

func GetCarbonIntensity(w http.ResponseWriter, now string, nowPlusIntervalInMinutes string) (CarbonintensityResponse, error) {
	utils.Logger().Info("GetCarbonIntensity", zap.Field{Key: "now", Type: zapcore.StringType, String: now}, zap.Field{Key: "nowPlusIntervalInMinutes", Type: zapcore.StringType, String: nowPlusIntervalInMinutes})
	result, err := http.Get(fmt.Sprintf("https://api.carbonintensity.org.uk/intensity/%s/%s", now, nowPlusIntervalInMinutes))

	if err != nil {
		utils.Logger().Error("Error getting carbon intensity", zap.Error(err))
		return CarbonintensityResponse{Data: nil}, err
	}

	defer result.Body.Close()

	body, err := ioutil.ReadAll(result.Body)

	if err != nil {
		utils.Logger().Error("Error reading carbon intensity response", zap.Error(err))
		return CarbonintensityResponse{Data: nil}, err
	}

	var data CarbonintensityResponse

	err = json.Unmarshal(body, &data)

	if err != nil {
		utils.Logger().Error("Error unmarshalling carbon intensity response", zap.Error(err))
		return CarbonintensityResponse{Data: nil}, err
	}

	utils.Logger().Info("Successfully got carbon intensity")
	return data, nil
}
