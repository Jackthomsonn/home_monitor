package functions

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"cloud.google.com/go/pubsub"
)

type Data struct {
	Temperature float32 `json:"temperature"`
	ClientId  string `json:"client_id"`
	Timestamp string `json:"timestamp"`
}

type ResponseType struct {
	Type string `json:"type"`
	Id string `json:"id"`
}

type ErrorResponseType struct {
	Type string `json:"type"`
	Message string `json:"message"`
}

func PublishData(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil { throwError(w, err); return }

	var data Data
	err = json.Unmarshal(body, &data)

	if err != nil { throwError(w, err); return }

	id, err := publishDataToGCP(body)

	if err != nil { throwError(w, err); return }

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ResponseType{Type: "Success", Id: id})

	return
}

func publishDataToGCP(body []byte) (string, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, "home-monitor-373013")

	if err != nil {
		return "", errors.New(err.Error())
	}

	topic := client.Topic("state")

	result := topic.Publish(ctx, &pubsub.Message{
		Data: body,
	})

	id, err := result.Get(context.Background())

	return id, nil
}

func throwError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ErrorResponseType{Message: err.Error(), Type: "failed"})
}