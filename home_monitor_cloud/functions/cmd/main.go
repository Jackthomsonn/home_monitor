package main

import (
	"log"
	"net/http"
	"os"

	"jackthomson.com/functions"
)

func main() {
	http.HandleFunc("/publishData", functions.PublishData)
	http.HandleFunc("/performCheck", functions.PerformCheck)
	http.HandleFunc("/ingestConsumptionData", functions.IngestConsumptionData)
	http.HandleFunc("/ingestCarbonIntensityData", functions.IngestCarbonIntensityData)
	http.HandleFunc("/getTotalsForHome", functions.GetTotalsForHome)
	http.HandleFunc("/triggerConsumptionData", functions.TriggerConsumptionData)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
