package utils

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func PublishDataToGCP(data interface{}, topicName string) (string, error) {
	Logger().Info("Publishing data to GCP", zap.Field{Key: "data", Type: zapcore.ReflectType, Interface: data}, zap.Field{Key: "topicName", Type: zapcore.ReflectType, Interface: topicName})

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, "home-monitor-373013")

	if err != nil {
		Logger().Error("Error creating pubsub client", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
		return "", err
	}

	defer client.Close()

	topic := client.Topic(topicName)
	dataBytes, err := json.Marshal(data)
	if err != nil {
		Logger().Error("Error marshalling data", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
		return "", err
	}

	result := topic.Publish(ctx, &pubsub.Message{
		Data: dataBytes,
	})

	id, err := result.Get(context.Background())
	if err != nil {
		Logger().Error("Error getting result from pubsub", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
		return "", err
	}

	Logger().Info("Published data to GCP", zap.Field{Key: "data", Type: zapcore.ReflectType, Interface: data}, zap.Field{Key: "topicName", Type: zapcore.ReflectType, Interface: topicName}, zap.Field{Key: "id", Type: zapcore.ReflectType, Interface: id})

	return id, nil
}
