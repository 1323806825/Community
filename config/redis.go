package config

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"time"
)

var rdClient *redis.Client
var nDuration = 30 * 24 * 60 * 60 * time.Second

type RedisClient struct {
}

func InitRedis() (*RedisClient, error) {
	rdClient = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("db.redis"),
		Password: "",
		DB:       0,
	})

	_, err := rdClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return &RedisClient{}, nil
}

func (m *RedisClient) GetKeysAndValue(pattern string) (map[string]string, error) {
	keys, err := rdClient.Keys(context.Background(), pattern).Result()
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, key := range keys {
		value, err := rdClient.Get(context.Background(), key).Result()
		if err != nil {
			return nil, err
		}
		result[key] = value
	}
	return result, nil

}

func (m *RedisClient) Delete(key ...string) error {
	//支持批量删除
	return rdClient.Del(context.Background(), key...).Err()
}

func (m *RedisClient) Set(key string, value any, rest ...any) error {
	d := nDuration
	if len(rest) > 0 {
		if v, ok := rest[0].(time.Duration); ok {
			d = v
		}
	}
	return rdClient.Set(context.Background(), key, value, d).Err()
}

func (m *RedisClient) Get(key string) (any, error) {
	return rdClient.Get(context.Background(), key).Result()
}
