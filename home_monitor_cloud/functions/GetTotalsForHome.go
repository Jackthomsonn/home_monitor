package functions

import (
	"encoding/json"
	"net/http"
	"os"

	"jackthomson.com/functions/services"
)

func GetTotalsForHome(w http.ResponseWriter, r *http.Request) {
	var originToUse string = "https://jackthomson.co.uk"

	if os.Getenv("DEVELOPMENT_MODE") == "true" {
		originToUse = "http://localhost:3000"
	}

	w.Header().Set("Access-Control-Allow-Origin", originToUse)
	var response, err = services.HomeTotals()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
