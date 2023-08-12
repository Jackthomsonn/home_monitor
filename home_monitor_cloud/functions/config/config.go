package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	StrategyToUse string
}

func ReadConfig() (Config, error) {
	jsonFile, err := os.Open("config/config.json")

	c := Config{}

	if err != nil {
		return Config{}, err
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &c)

	if err != nil {
		return Config{}, err
	}

	return c, nil
}
