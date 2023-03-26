package functions

import (
	"encoding/json"
	"net/http"

	"jackthomson.com/functions/models"
	"jackthomson.com/functions/services"
)

func IngestHomeTotals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()

	key, err := services.IngestHomeTotals()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Error{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseType{Type: "Success", Id: key.String()})
}
