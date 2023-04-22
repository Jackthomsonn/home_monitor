package utils

import (
	"context"

	"cloud.google.com/go/bigquery"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InsertDataIntoBiqQuery(ctx context.Context, data any, tableId string) error {
	client, err := bigquery.NewClient(ctx, "home-monitor-373013")

	if err != nil {
		Logger().Error("Error creating bigquery client", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
		return err
	}

	defer client.Close()

	table := client.Dataset("home_monitor_dataset").Table(tableId)

	if insertErr := table.Inserter().Put(ctx, data); insertErr != nil {
		Logger().Error("Error inserting data into bigquery", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: insertErr})
		return insertErr
	}

	return nil
}
