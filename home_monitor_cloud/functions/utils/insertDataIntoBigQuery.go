package utils

import (
	"context"

	"cloud.google.com/go/bigquery"
)

func InsertDataIntoBiqQuery(ctx context.Context, data any, tableId string) error {
  client, err := bigquery.NewClient(ctx, "home-monitor-373013")
  
  if err != nil { return err }

  defer client.Close()

  table := client.Dataset("home_monitor_dataset").Table(tableId)

  if insertErr := table.Inserter().Put(ctx, data); insertErr != nil {
    return insertErr
  }

  return nil
}