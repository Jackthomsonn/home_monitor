package functions

import (
	"encoding/json"
	"net/http"
	"os"

	"jackthomson.com/functions/models"
	"jackthomson.com/functions/services"
)

func GetTotalsForHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var originToUse string = "https://jackthomson.co.uk"

	if os.Getenv("DEVELOPMENT_MODE") == "true" {
		originToUse = "http://localhost:3000"
	}

	w.Header().Set("Access-Control-Allow-Origin", originToUse)
	response, err := services.HomeTotals()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Error{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
