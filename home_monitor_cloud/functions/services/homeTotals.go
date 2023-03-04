package services

import (
	"context"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

type UserTotalsResponse struct {
	ConsumptionTotal float64 `json:"consumptionTotal,omitempty"`
	CarbonTotal      float64 `json:"carbonTotal,omitempty"`
}

type RowResponse struct {
	CarbonIntensity float64 `json:"carbonIntensity,omitempty"`
	Consumption     float64 `json:"consumption,omitempty"`
}

// This will turn into an ETL that will store the data in Postgres
func HomeTotals() (UserTotalsResponse, error) {
	client, err := bigquery.NewClient(context.Background(), "home-monitor-373013")

	if err != nil {
		return UserTotalsResponse{}, err
	}

	defer client.Close()

	query := client.Query("SELECT `home-monitor-373013.home_monitor_dataset.home_monitor_carbon_intensity`.timestamp as ts, MAX(`home-monitor-373013.home_monitor_dataset.home_monitor_carbon_intensity`.actual) as carbonIntensity, MAX(`home_monitor_dataset.home_monitor_consumption_table`.value) as consumption, FROM `home-monitor-373013.home_monitor_dataset.home_monitor_carbon_intensity` INNER JOIN `home_monitor_dataset.home_monitor_consumption_table` ON `home-monitor-373013.home_monitor_dataset.home_monitor_carbon_intensity`.timestamp = `home_monitor_dataset.home_monitor_consumption_table`.timestamp WHERE DATE(`home-monitor-373013.home_monitor_dataset.home_monitor_carbon_intensity`.timestamp) = DATE_SUB(CURRENT_DATE(), INTERVAL 1 DAY) AND Date(`home-monitor-373013.home_monitor_dataset.home_monitor_carbon_intensity`.timestamp) >= DATE_TRUNC(DATE_SUB(CURRENT_DATE(), INTERVAL 1 DAY), DAY) GROUP BY ts ORDER BY ts ASC")

	it, err := query.Read(context.Background())

	if err != nil {
		return UserTotalsResponse{}, err
	}

	var carbonTotal float64
	var consumptionTotal float64

	for {
		var row RowResponse
		err := it.Next(&row)

		if err == iterator.Done {
			break
		}

		if err != nil {
			return UserTotalsResponse{}, err
		}

		carbonTotal += row.CarbonIntensity * row.Consumption
		consumptionTotal += row.Consumption
	}

	return UserTotalsResponse{CarbonTotal: carbonTotal, ConsumptionTotal: consumptionTotal}, nil
}
