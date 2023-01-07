package utils

import (
	"context"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

func GetSecret(secret_name string, ctx context.Context) (response string, error error) {
  secret_manager_client, err := secretmanager.NewClient(ctx)

  if err != nil {
    return "", err
  }
  
  defer secret_manager_client.Close()

  secret_req := &secretmanagerpb.AccessSecretVersionRequest{
    Name: secret_name,
  }

  result, err := secret_manager_client.AccessSecretVersion(ctx, secret_req)

  if err != nil {
    return "", err
  }

  return string(result.Payload.Data), nil
}