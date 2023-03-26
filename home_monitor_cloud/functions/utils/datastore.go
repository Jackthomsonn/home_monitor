package utils

import (
	"context"

	"cloud.google.com/go/datastore"
)

type Total struct {
	ConsumptionTotal float64 `json:"consumptionTotal,omitempty"`
	CarbonTotal      float64 `json:"carbonTotal,omitempty"`
}

func WriteToDatastore(name_key *datastore.Key, data interface{}) (datastore.Key, error) {
	data_store, data_store_err := datastore.NewClient(context.Background(), "home-monitor-373013")

	if data_store_err != nil {
		return datastore.Key{}, data_store_err
	}

	key, err := data_store.Put(context.Background(), name_key, data)

	if err != nil {
		return datastore.Key{}, err
	}

	return *key, nil
}

func ReadFromDatastore(name_key *datastore.Key) (interface{}, error) {
	data_store, data_store_err := datastore.NewClient(context.Background(), "home-monitor-373013")

	if data_store_err != nil {
		return nil, data_store_err
	}

	var data Total = Total{CarbonTotal: 0, ConsumptionTotal: 0}

	err := data_store.Get(context.Background(), name_key, &data)

	if err != nil {
		return nil, err
	}

	return data, nil
}
