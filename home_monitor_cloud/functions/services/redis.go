package services

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"jackthomson.com/functions/utils"
)

func GetDataFromRedis(dest interface{}, key string) error {
	ctx := context.Background()
	redisConnectionString, err := utils.GetSecret("projects/345305797254/secrets/redis_connection_string/versions/latest", ctx)

	if err != nil {
		return err
	}

	opt, _ := redis.ParseURL(redisConnectionString)
	client := redis.NewClient(opt)

	val := client.Get(ctx, key)

	if val.Val() == "" {
		return errors.New("no data found")
	}

	if err := json.Unmarshal([]byte(val.Val()), &dest); err != nil {
		return err
	}

	return nil
}

func CreateDataInRedis(data interface{}, key string, ttl time.Duration) error {
	utils.Logger().Info("Creating data in redis", zap.Field{Key: "key", Type: zapcore.StringType, String: key}, zap.Field{Key: "ttl", Type: zapcore.Int64Type, Integer: int64(ttl.Hours())})

	ctx := context.Background()
	redisConnectionString, err := utils.GetSecret("projects/345305797254/secrets/redis_connection_string/versions/latest", ctx)

	if err != nil {
		utils.Logger().Error("Error getting redis connection string", zap.Error(err))
		return err
	}

	opt, _ := redis.ParseURL(redisConnectionString)
	client := redis.NewClient(opt)

	jsonData, err := json.Marshal(data)

	if err != nil {
		utils.Logger().Error("Error marshalling data", zap.Error(err))
		return err
	}

	if err := client.Set(ctx, key, jsonData, ttl).Err(); err != nil {
		utils.Logger().Error("Error setting data in redis", zap.Error(err))
		return err
	}

	utils.Logger().Info("Created data in redis", zap.Field{Key: "key", Type: zapcore.StringType, String: key}, zap.Field{Key: "ttl", Type: zapcore.Int64Type, Integer: int64(ttl.Hours())})

	return nil
}

func RemoveDataFromRedis(key string) error {
	utils.Logger().Info("Removing data from redis", zap.Field{Key: "key", Type: zapcore.StringType, String: key})

	ctx := context.Background()
	redisConnectionString, err := utils.GetSecret("projects/345305797254/secrets/redis_connection_string/versions/latest", ctx)

	if err != nil {
		utils.Logger().Error("Error getting redis connection string", zap.Error(err))
		return err
	}

	opt, _ := redis.ParseURL(redisConnectionString)
	client := redis.NewClient(opt)

	if err := client.Del(ctx, key).Err(); err != nil {
		utils.Logger().Error("Error removing data from redis", zap.Error(err))
		return err
	}

	utils.Logger().Info("Removed data from redis", zap.Field{Key: "key", Type: zapcore.StringType, String: key})

	return nil
}
