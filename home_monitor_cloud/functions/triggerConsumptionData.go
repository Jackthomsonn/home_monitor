package functions

import (
	"encoding/json"
	"net/http"

	"jackthomson.com/functions/models"
	"jackthomson.com/functions/utils"
)

func TriggerConsumptionData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data Data
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Error{Message: err.Error()})
		return
	}

	defer r.Body.Close()

	id, err := utils.PublishDataToGCP(nil, "consumption-ingestion")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Error{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseType{Type: "Success", Id: id})
}
