package functions

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"jackthomson.com/functions/models"
	"jackthomson.com/functions/utils"
)

type ConsumptionValues struct {
	Timestamp string  `json:"timestamp"`
	Value     float32 `json:"value"`
}

type AvailableCacheRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type ConsumptionResponse struct {
	Unit                string              `json:"unit"`
	Granularity         string              `json:"granularity"`
	Start               string              `json:"start"`
	End                 string              `json:"end"`
	ResponseTimestamp   string              `json:"responseTimestamp"`
	Resource            string              `json:"resource"`
	Values              []ConsumptionValues `json:"values"`
	AvailableCacheRange AvailableCacheRange `json:"availableCacheRange"`
}

func IngestConsumptionData(w http.ResponseWriter, r *http.Request) {
	utils.Logger().Info("IngestConsumptionData", zap.Field{Key: "method", Type: zapcore.StringType, String: r.Method}, zap.Field{Key: "url", Type: zapcore.StringType, String: r.URL.String()})
	w.Header().Set("Content-Type", "application/json")

	response, err := getConsumptionData(w)

	if err != nil {
		utils.Logger().Error("Error getting consumption data", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Error{Message: err.Error()})
		return
	}

	response.Values = response.Values[:len(response.Values)-1]

	bigqueryErr := utils.InsertDataIntoBiqQuery(context.Background(), response.Values, "home_monitor_consumption_table")

	if bigqueryErr != nil {
		utils.Logger().Error("Error inserting consumption data", zap.Error(bigqueryErr))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Error{Message: err.Error()})
		return
	}

	utils.Logger().Info("Successfully ingested consumption data", zap.Field{Key: "start", Type: zapcore.StringType, String: response.Start}, zap.Field{Key: "end", Type: zapcore.StringType, String: response.End})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func getStartAndEndDates(daysFromNow int) (string, string) {
	utils.Logger().Info("Getting start and end dates", zap.Field{Key: "daysFromNow", Type: zapcore.Int64Type, Integer: int64(daysFromNow)})
	currentTime := time.Now()

	start := currentTime.AddDate(0, 0, -1).UTC().Format("20060102")

	end := currentTime.UTC().Format("20060102")

	utils.Logger().Info("Successfully got start and end dates", zap.Field{Key: "start", Type: zapcore.StringType, String: start}, zap.Field{Key: "end", Type: zapcore.StringType, String: end})

	return start, end
}

func getConsumptionData(w http.ResponseWriter) (ConsumptionResponse, error) {
	start, end := getStartAndEndDates(1)

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://consumer-api.data.n3rgy.com/electricity/consumption/1?start="+start+"&end="+end+"&output=json", nil)

	if err != nil {
		utils.Logger().Error("Error creating request", zap.Error(err))
		return ConsumptionResponse{}, err
	}

	secret, err := utils.GetSecret("projects/345305797254/secrets/consumption_secret/versions/latest", context.Background())

	if err != nil {
		utils.Logger().Error("Error getting secret", zap.Error(err))
		return ConsumptionResponse{}, err
	}

	req.Header.Set("Authorization", secret)

	res, err := client.Do(req)

	if err != nil {
		utils.Logger().Error("Error making request", zap.Error(err))
		return ConsumptionResponse{}, err
	}

	defer client.CloseIdleConnections()

	defer res.Body.Close()

	response := ConsumptionResponse{}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		utils.Logger().Error("Error reading response body", zap.Error(err))
		return ConsumptionResponse{}, err
	}

	err = json.Unmarshal(body, &response)

	if err != nil {
		utils.Logger().Error("Error unmarshalling response body", zap.Error(err))
		return ConsumptionResponse{}, err
	}

	if len(response.Values) == 0 {
		utils.Logger().Error("No consumption data found", zap.Field{Key: "start", Type: zapcore.StringType, String: response.Start}, zap.Field{Key: "end", Type: zapcore.StringType, String: response.End})
		return ConsumptionResponse{}, errors.New("downstream error: no consumption data found")
	}

	utils.Logger().Info("Successfully retrieved consumption data", zap.Field{Key: "start", Type: zapcore.StringType, String: response.Start}, zap.Field{Key: "end", Type: zapcore.StringType, String: response.End})

	if err := utils.RemoveDataFromRedis("Total"); err != nil {
		utils.Logger().Error("Error removing data from redis", zap.Error(err))
	}

	return response, nil
}
