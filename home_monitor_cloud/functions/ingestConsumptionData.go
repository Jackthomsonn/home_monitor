package functions

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"cloud.google.com/go/bigquery"
	"jackthomson.com/functions/utils"
)

type ConsumptionValues struct {
  Timestamp string `json:"timestamp"`
  Value    float32    `json:"value"`
}

type AvailableCacheRange struct {
  Start string `json:"start"`
  End string `json:"end"`
}

type ConsumptionResponse struct {
  Unit    string `json:"unit"`
  Granularity   string `json:"granularity"`
  Start  string `json:"start"`
  End string `json:"end"`
  ResponseTimestamp string `json:"responseTimestamp"`
  Resource string `json:"resource"`
  Values []ConsumptionValues `json:"values"`
  AvailableCacheRange AvailableCacheRange `json:"availableCacheRange"`
}

func IngestConsumptionData(w http.ResponseWriter, r *http.Request) {
  err, response := getConsumptionData(w)

  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  bigqueryErr := insertDataIntoBiqQuery(context.Background(), response.Values)

  if bigqueryErr != nil {
    http.Error(w, bigqueryErr.Error(), http.StatusBadRequest)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(response)
}

func getStartAndEndDates(daysFromNow int) (string, string) {
  currentTime := time.Now()

  start := currentTime.UTC().Format("20060102")
  
  end := currentTime.AddDate(0, 0, daysFromNow).UTC().Format("20060102")

  return start, end
}

func getConsumptionData(w http.ResponseWriter) (error, ConsumptionResponse) {
  start, end := getStartAndEndDates(1)

  client := &http.Client{}
  req, err := http.NewRequest("GET", "https://consumer-api.data.n3rgy.com/electricity/consumption/1?start=" + start + "&end=" + end + "&output=json", nil)

  if err != nil {
    return err, ConsumptionResponse{}
  }

  secret, err := utils.GetSecret("projects/345305797254/secrets/consumption_secret/versions/latest", context.Background())

  if err != nil {
    return err, ConsumptionResponse{}
  }

  req.Header.Set("Authorization", secret)

  res, err := client.Do(req)

  if err != nil {
    return err, ConsumptionResponse{}
  }

  defer res.Body.Close()

  response := ConsumptionResponse{}

  body, err := ioutil.ReadAll(res.Body)

  if err != nil {
    return err, ConsumptionResponse{}
  }

  err = json.Unmarshal(body, &response)

  if err != nil {
    return err, ConsumptionResponse{}
  }

  return nil, response
}

func insertDataIntoBiqQuery(ctx context.Context, consumptionData []ConsumptionValues) error {
  client, err := bigquery.NewClient(ctx, "home-monitor-373013")
  
  if err != nil { return err }

  table := client.Dataset("home_monitor_dataset").Table("home_monitor_consumption_table")

  insertErr := table.Inserter().Put(ctx, consumptionData)

  if insertErr != nil { return insertErr }

  return nil
}