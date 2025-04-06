package rediscl

import (
	"github.com/redis/go-redis/v9"
	"github.com/ruziba3vich/itv_test_project/pkg/config"
)

func NewRedisClient(cfg *config.Config) *redis.Client {
	// Redis client
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
}
