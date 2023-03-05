package functions

import (
	"context"
	"encoding/json"
	"net/http"

	"cloud.google.com/go/pubsub"
	"jackthomson.com/functions/models"
)

type Data struct {
	Temperature *float32 `json:"temperature"`
	ClientId    *string  `json:"client_id"`
	Timestamp   *string  `json:"timestamp"`
}

type ResponseType struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

func PublishData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data Data
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Error{Message: err.Error()})
		return
	}

	defer r.Body.Close()

	id, err := publishDataToGCP(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Error{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseType{Type: "Success", Id: id})
}

func publishDataToGCP(data Data) (string, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, "home-monitor-373013")
	if err != nil {
		return "", err
	}

	defer client.Close()

	topic := client.Topic("state")
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
