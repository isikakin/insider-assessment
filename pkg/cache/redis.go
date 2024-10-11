package cache

import (
	"assesment/pkg/config"
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type Cache interface {
	Set(k string, x interface{}, d time.Duration)
	Retrieve(k string) (interface{}, bool)
	Exists(k string) bool
}

type redisDefault struct {
	client *redis.Client
	ctx    context.Context
}

func (c *redisDefault) Set(k string, x interface{}, d time.Duration) {
	c.client.Set(context.Background(), k, x, d)
}

func (c *redisDefault) Retrieve(k string) (interface{}, bool) {
	val, err := c.client.Get(c.ctx, k).Result()
	if errors.Is(err, redis.Nil) {
		return nil, false
	} else if err != nil {
		return nil, false
	}

	return val, true
}

func (c *redisDefault) Exists(k string) bool {
	exists, err := c.client.Exists(c.ctx, k).Result()
	if errors.Is(err, redis.Nil) {
		return false
	}

	if exists == 1 {
		return true
	}
	return false
}

func NewDistributed(configuration config.Configuration) Cache {

	return &redisDefault{
		client: redis.NewClient(&redis.Options{
			Addr: configuration.Redis.Url,
		}),
		ctx: context.Background(),
	}
}
