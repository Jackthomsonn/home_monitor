package utils

import (
	"context"

	"cloud.google.com/go/datastore"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type DataStoreDevice struct {
	A          int
	K          *datastore.Key `datastore:"__key__"`
	Ip         []int          `json:"ip"`
	Alias      string         `json:"alias"`
	Feature    string         `json:"feature"`
	OnTime     int            `json:"on_time"`
	RelayState int            `json:"relay_state"`
	DeviceId   string         `json:"device_id"`
	DeviceType string         `json:"device_type"`
}

func WriteToDatastore(name_key *datastore.Key, data interface{}) (datastore.Key, error) {
	Logger().Info("Writing to datastore", zap.Field{Key: "name_key", Type: zapcore.ReflectType, Interface: name_key}, zap.Field{Key: "data", Type: zapcore.ReflectType, Interface: data})

	data_store, data_store_err := datastore.NewClient(context.Background(), "home-monitor-373013")

	if data_store_err != nil {
		Logger().Error("Error creating datastore client", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: data_store_err})
		return datastore.Key{}, data_store_err
	}

	key, err := data_store.Put(context.Background(), name_key, data)

	if err != nil {
		Logger().Error("Error writing to datastore", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
		return datastore.Key{}, err
	}

	return *key, nil
}

func ReadFromDatastore(nameKey *datastore.Key, dest interface{}) error {
	Logger().Info("Reading from datastore", zap.Field{Key: "nameKey", Type: zapcore.ReflectType, Interface: nameKey})

	client, err := datastore.NewClient(context.Background(), "home-monitor-373013")

	if err != nil {
		Logger().Error("Error creating datastore client", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
		return err
	}

	if err := client.Get(context.Background(), nameKey, dest); err != nil {
		Logger().Error("Error getting data from datastore", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
		return err
	}

	return nil
}

func ReadAllFromDataStore() ([]DataStoreDevice, error) {
	Logger().Info("Reading all from datastore")

	client, err := datastore.NewClient(context.Background(), "home-monitor-373013")

	if err != nil {
		Logger().Error("Error creating datastore client", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
		return nil, err
	}

	var entities []DataStoreDevice

	_, err = client.GetAll(context.Background(), datastore.NewQuery("Device"), &entities)

	if err != nil {
		Logger().Error("Error getting data from datastore", zap.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
		return nil, err
	}

	return entities, nil
}
