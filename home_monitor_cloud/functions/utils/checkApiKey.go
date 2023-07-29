package utils

import (
	"context"
	"errors"
)

func CheckApiKey(api_key string) error {
	secret_api_key, err := GetSecret("projects/345305797254/secrets/api_key/versions/latest", context.Background())

	if err != nil {
		return err
	}

	if api_key != secret_api_key {
		return errors.New("Invalid API key")
	}
}
