package services

import (
	"time"

	"cloud.google.com/go/datastore"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"jackthomson.com/functions/models"
	"jackthomson.com/functions/utils"
)

func GetTotalsForHome() (interface{}, error) {
	utils.Logger().Info("GetTotalsForHome")

	var key = datastore.NameKey("Total", "total", nil)
	var homeTotals models.HomeTotal = models.HomeTotal{CarbonTotal: 0, ConsumptionTotal: 0}

	if err := utils.GetDataFromRedis(&homeTotals, "Total"); err == nil {
		utils.Logger().Info("Returning totals from cache", zap.Field{Key: "data", Type: zapcore.ReflectType, Interface: homeTotals})
		return homeTotals, err
	}

	if err := utils.ReadFromDatastore(key, &homeTotals); err != nil {
		utils.Logger().Error("Error getting totals for home", zap.Error(err))
		return models.HomeTotal{}, err
	}

	if err := storeTotalsForHomeInCache(&homeTotals); err != nil {
		utils.Logger().Error("Error storing totals for home in cache", zap.Error(err))
		return models.HomeTotal{}, err
	}

	return homeTotals, nil
}

func storeTotalsForHomeInCache(data *models.HomeTotal) error {
	utils.Logger().Info("Storing totals for home in cache", zap.Field{Key: "data", Type: zapcore.ReflectType, Interface: data})

	delta := time.Duration(5) * time.Hour

	if err := utils.CreateDataInRedis(data, "Total", delta); err != nil {
		utils.Logger().Error("Error storing totals for home in cache", zap.Error(err))
		return err
	}

	utils.Logger().Info("Successfully stored totals for home in cache", zap.Field{Key: "delta", Type: zapcore.Int64Type, Integer: int64(delta.Hours())})

	return nil
}
