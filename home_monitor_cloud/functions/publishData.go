package functions

import (
	"context"
	"encoding/json"
	"net/http"

	"cloud.google.com/go/pubsub"
)

type Data struct {
	Temperature *float32 `json:"temperature"`
	ClientId  *string `json:"client_id"`
	Timestamp *string `json:"timestamp"`
}

type ResponseType struct {
	Type string `json:"type"`
	Id string `json:"id"`
}

func PublishData(w http.ResponseWriter, r *http.Request) {
	var data Data
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	id, err := publishDataToGCP(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
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