package redisx

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tOnkowzl/libs/logx"
)

type Client struct {
	client *redis.Client
}

func NewClient(opt *redis.Options) *Client {
	return &Client{
		client: redis.NewClient(opt),
	}
}

func (c *Client) Get(ctx context.Context, key string) *Bind {
	start := time.Now()
	res, err := c.client.Get(ctx, key).Result()

	logx.WithContext(ctx).WithFields(logrus.Fields{
		"key":      key,
		"value":    res,
		"duration": time.Since(start).String(),
	}).Info("redis get information")

	return &Bind{
		Result: res,
		Err:    errors.WithStack(err),
	}
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	start := time.Now()
	err := c.client.Set(ctx, key, value, expiration).Err()

	logx.WithContext(ctx).WithFields(logrus.Fields{
		"key":      key,
		"value":    value,
		"duration": time.Since(start).String(),
	}).Info("redis set information")

	return errors.WithStack(err)
}

type Bind struct {
	Result string
	Err    error
}

func (b *Bind) Bind(i interface{}) error {
	if b.Err != nil {
		return b.Err
	}

	err := json.Unmarshal([]byte(b.Result), i)
	return errors.WithStack(err)
}
