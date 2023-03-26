package services

import (
	"fmt"

	"cloud.google.com/go/datastore"
	"jackthomson.com/functions/utils"
)

func GetTotalsForHome() (interface{}, error) {
	key := datastore.NameKey("Total", "total", nil)
	fmt.Println(key)
	fmt.Println(*key)
	fmt.Println(&key)
	data, err := utils.ReadFromDatastore(key)

	if err != nil {
		return datastore.Key{}, err
	}

	return data, nil
}
