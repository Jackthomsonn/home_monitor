package utils

import (
	"context"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"go.uber.org/zap/zapcore"
)

func GetSecret(secret_name string, ctx context.Context) (string, error) {
	Logger().Info("Getting secret", zapcore.Field{Key: "secret_name", Type: zapcore.StringType, String: secret_name}, zapcore.Field{Key: "ctx", Type: zapcore.ReflectType, Interface: ctx})

	secret_manager_client, err := secretmanager.NewClient(ctx)

	if err != nil {
		Logger().Error("Error creating secret manager client", zapcore.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
		return "", err
	}

	defer secret_manager_client.Close()

	secret_req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secret_name,
	}

	result, err := secret_manager_client.AccessSecretVersion(ctx, secret_req)

	if err != nil {
		Logger().Error("Error getting secret", zapcore.Field{Key: "error", Type: zapcore.ReflectType, Interface: err})
		return "", err
	}

	Logger().Info("Successfully got secret", zapcore.Field{Key: "secret_name", Type: zapcore.StringType, String: secret_name}, zapcore.Field{Key: "ctx", Type: zapcore.ReflectType, Interface: ctx})

	return string(result.Payload.Data), nil
}
