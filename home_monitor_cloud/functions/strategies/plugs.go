package strategies

import (
	"context"
	"time"

	"go.uber.org/zap"
	"jackthomson.com/functions/models"
	"jackthomson.com/functions/services"
	"jackthomson.com/functions/utils"
)

func GetPlugConsumptionData() (models.ConsumptionResponse, error) {
	query := `SELECT REGEXP_EXTRACT(payload, r'alias: ([^,}]+)') AS alias, REGEXP_EXTRACT(payload, r'power_mw: (\d+)') AS power_mw, timestamp FROM ` + "`home-monitor-373013.home_monitor_dataset.home_monitor_table`" + ` WHERE timestamp >= TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 24 HOUR) ORDER BY timestamp DESC`

	response, err := services.GetDataFromBigQuery[models.EnergyConsumption](context.Background(), "home_monitor_table", query)

	if err != nil {
		utils.Logger().Error("Error getting consumption data", zap.Error(err))
		return models.ConsumptionResponse{}, err
	}

	values := []models.ConsumptionValues{}

	for i, v := range response {
		if i == 0 {
			continue
		}

		values = append(values, models.ConsumptionValues{
			Value:     v.PowerMw,
			Timestamp: v.Timestamp.Format(time.RFC3339),
		})
	}

	return models.ConsumptionResponse{
		Values:            values,
		Start:             response[len(response)-1].Timestamp.Format(time.RFC3339),
		End:               response[0].Timestamp.Format(time.RFC3339),
		Granularity:       "oneMinute",
		ResponseTimestamp: time.Now().Format(time.RFC3339),
		Resource:          "plug",
		Unit:              "mW",
	}, nil
}
