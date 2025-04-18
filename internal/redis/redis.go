package redis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"time"
	"urlShortener/internal/config"
)

type Client struct {
	cfg    *config.RedisSettings
	client *redis.Client
}

func New(cfg *config.RedisSettings) *Client {
	cl := Client{
		cfg: cfg,
	}

	cl.client = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: "",
		DB:       0,
	})

	return &cl
}

func (c *Client) Get(key string) ([]byte, error) {
	cmd := c.client.Get(context.Background(), key)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	return []byte(cmd.Val()), nil
}

func (c *Client) Set(key string, v any, time time.Duration) error {
	val, err := json.Marshal(v)
	if err != nil {
		return err
	}

	err = c.client.Set(context.Background(), key, string(val), time).Err()
	if err != nil {
		return err
	}

	return nil
}
