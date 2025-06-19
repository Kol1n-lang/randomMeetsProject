package database

import (
	"github.com/redis/go-redis/v9"
	"randomMeetsProject/config"
)

func RedisClient() (*redis.Client, error) {
	cfg, err := config.LoadConfig("config.toml")
	if err != nil {
		return nil, err
	}
	client := cfg.RedisClient()
	return client, nil
}
