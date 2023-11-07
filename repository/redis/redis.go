package redis

import (
	"github.com/go-redis/redis/v8"
)

func New(cfg *Cfg) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

}

type Cfg struct {
	Address  string
	Password string
	DB       int
	PoolSize int
}
