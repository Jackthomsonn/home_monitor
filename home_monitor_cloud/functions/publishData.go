package functions

import (
	"encoding/json"
	"net/http"

	"jackthomson.com/functions/models"
	"jackthomson.com/functions/utils"
)

type Data struct {
	Payload   *string `json:"payload"`
	Topic     *string `json:"topic"`
	Type      *string `json:"type"`
	ClientId  *string `json:"client_id"`
	Timestamp *string `json:"timestamp"`
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

	id, err := utils.PublishDataToGCP(data, "state")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Error{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseType{Type: "Success", Id: id})
}
