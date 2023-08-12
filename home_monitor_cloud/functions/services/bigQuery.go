package services

import (
	"context"

	"cloud.google.com/go/bigquery"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/api/iterator"
	"jackthomson.com/functions/utils"
)

func InsertDataIntoBiqQuery(ctx context.Context, data any, tableId string) error {
	client, err := bigquery.NewClient(ctx, "home-monitor-373013")

	if err != nil {
		utils.Logger().Error("Error creating bigquery client", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
		return err
	}

	defer client.Close()

	table := client.Dataset("home_monitor_dataset").Table(tableId)

	if insertErr := table.Inserter().Put(ctx, data); insertErr != nil {
		utils.Logger().Error("Error inserting data into bigquery", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: insertErr})
		return insertErr
	}

	return nil
}

func GetDataFromBigQuery[T any](ctx context.Context, tableId string, q string) ([]T, error) {
	client, err := bigquery.NewClient(ctx, "home-monitor-373013")

	var row []T

	if err != nil {
		utils.Logger().Error("Error creating bigquery client", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
		return row, err
	}

	defer client.Close()

	query := client.Query(q)

	it, err := query.Read(context.Background())

	if err != nil {
		utils.Logger().Error("Error reading bigquery data", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
	}

	for {
		var nRow T
		err := it.Next(&nRow)
		if err == iterator.Done {
			break
		}
		if err != nil {
			utils.Logger().Info(err.Error())
			utils.Logger().Error("Error iterating bigquery data", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
		}
		row = append(row, nRow)
	}

	return row, nil
}
