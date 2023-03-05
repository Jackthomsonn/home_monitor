package utils

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
)

func PublishDataToGCP(data interface{}, topicName string) (string, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, "home-monitor-373013")
	if err != nil {
		return "", err
	}

	defer client.Close()

	topic := client.Topic(topicName)
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	result := topic.Publish(ctx, &pubsub.Message{
		Data: dataBytes,
	})

	id, err := result.Get(context.Background())
	if err != nil {
		return "", err
	}

	return id, nil
}
