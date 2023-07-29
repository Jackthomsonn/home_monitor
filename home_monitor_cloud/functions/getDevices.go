package functions

import (
	"encoding/json"
	"log"
	"net/http"

	"jackthomson.com/functions/utils"
)

func GetDevices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Println("Getting devices")
	keys, err := utils.ReadAllFromDataStore()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(keys)
}
