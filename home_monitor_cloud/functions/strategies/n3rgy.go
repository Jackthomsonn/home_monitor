package strategies

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"jackthomson.com/functions/models"
	"jackthomson.com/functions/utils"
)

func getStartAndEndDates(daysFromNow int) (string, string) {
	utils.Logger().Info("Getting start and end dates", zap.Field{Key: "daysFromNow", Type: zapcore.Int64Type, Integer: int64(daysFromNow)})
	currentTime := time.Now()

	start := currentTime.AddDate(0, 0, -1).UTC().Format("20060102")

	end := currentTime.UTC().Format("20060102")

	utils.Logger().Info("Successfully got start and end dates", zap.Field{Key: "start", Type: zapcore.StringType, String: start}, zap.Field{Key: "end", Type: zapcore.StringType, String: end})

	return start, end
}

func GetN3rgyConsumptionData() (models.ConsumptionResponse, error) {
	start, end := getStartAndEndDates(1)

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://consumer-api.data.n3rgy.com/electricity/consumption/1?start="+start+"&end="+end+"&output=json", nil)

	if err != nil {
		utils.Logger().Error("Error creating request", zap.Error(err))
		return models.ConsumptionResponse{}, err
	}

	secret, err := utils.GetSecret("projects/345305797254/secrets/consumption_secret/versions/latest", context.Background())

	if err != nil {
		utils.Logger().Error("Error getting secret", zap.Error(err))
		return models.ConsumptionResponse{}, err
	}

	req.Header.Set("Authorization", secret)

	res, err := client.Do(req)

	if err != nil {
		utils.Logger().Error("Error making request", zap.Error(err))
		return models.ConsumptionResponse{}, err
	}

	defer client.CloseIdleConnections()

	defer res.Body.Close()

	response := models.ConsumptionResponse{}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		utils.Logger().Error("Error reading response body", zap.Error(err))
		return models.ConsumptionResponse{}, err
	}

	err = json.Unmarshal(body, &response)

	if err != nil {
		utils.Logger().Error("Error unmarshalling response body", zap.Error(err))
		return models.ConsumptionResponse{}, err
	}

	if len(response.Values) == 0 {
		utils.Logger().Error("No consumption data found", zap.Field{Key: "start", Type: zapcore.StringType, String: response.Start}, zap.Field{Key: "end", Type: zapcore.StringType, String: response.End})
		return models.ConsumptionResponse{}, errors.New("downstream error: no consumption data found")
	}

	utils.Logger().Info("Successfully retrieved consumption data", zap.Field{Key: "start", Type: zapcore.StringType, String: response.Start}, zap.Field{Key: "end", Type: zapcore.StringType, String: response.End})

	return response, nil
}
