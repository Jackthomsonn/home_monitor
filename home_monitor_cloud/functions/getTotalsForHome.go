package functions

import (
	"encoding/json"
	"net/http"
	"time"

	"cloud.google.com/go/datastore"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"jackthomson.com/functions/models"
	"jackthomson.com/functions/services"
	"jackthomson.com/functions/utils"
)

func GetTotalsForHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	response, err := getTotalsForHome()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Error{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func getTotalsForHome() (interface{}, error) {
	utils.Logger().Info("GetTotalsForHome")

	var key = datastore.NameKey("Total", "total", nil)
	var homeTotals models.HomeTotal = models.HomeTotal{CarbonTotal: 0, ConsumptionTotal: 0}

	if err := services.GetDataFromRedis(&homeTotals, "Total"); err == nil {
		utils.Logger().Info("Returning totals from cache", zap.Field{Key: "data", Type: zapcore.ReflectType, Interface: homeTotals})
		return homeTotals, err
	}

	if err := services.ReadFromDatastore(key, &homeTotals); err != nil {
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

	if err := services.CreateDataInRedis(data, "Total", delta); err != nil {
		utils.Logger().Error("Error storing totals for home in cache", zap.Error(err))
		return err
	}

	utils.Logger().Info("Successfully stored totals for home in cache", zap.Field{Key: "delta", Type: zapcore.Int64Type, Integer: int64(delta.Hours())})

	return nil
}
