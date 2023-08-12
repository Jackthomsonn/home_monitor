package strategies

import "jackthomson.com/functions/models"

func StrategyFactory(strategy string) (models.ConsumptionResponse, error) {
	switch strategy {
	case "plug":
		return GetPlugConsumptionData()
	default:
		return GetN3rgyConsumptionData()
	}
}
