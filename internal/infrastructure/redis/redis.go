package redis

import (
	"doctor-go/internal/config"

	goredis "github.com/redis/go-redis/v9"
)

type Client struct {
	*goredis.Client
}

func New(cfg config.RedisConfig) *Client {
	client := goredis.NewClient(&goredis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	return &Client{Client: client}
}
