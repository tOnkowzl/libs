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
	*redis.Client
}

func NewClient(opt *redis.Options) *Client {
	return &Client{
		redis.NewClient(opt),
	}
}

func (c *Client) Get(ctx context.Context, key string) *Bind {
	start := time.Now()
	val, err := c.Client.Get(ctx, key).Result()

	logx.WithContext(ctx).WithFields(logrus.Fields{
		"key":      key,
		"value":    val,
		"duration": time.Since(start).String(),
	}).Info("redis get information")

	return &Bind{
		Val: val,
		Err: errors.WithStack(err),
	}
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	start := time.Now()
	err := c.Client.Set(ctx, key, value, expiration).Err()

	logx.WithContext(ctx).WithFields(logrus.Fields{
		"key":      key,
		"value":    value,
		"duration": time.Since(start).String(),
	}).Info("redis set information")

	return errors.WithStack(err)
}

func (c *Client) Del(ctx context.Context, keys ...string) error {
	start := time.Now()
	err := c.Client.Del(ctx, keys...).Err()

	logx.WithContext(ctx).WithFields(logrus.Fields{
		"keys":     keys,
		"duration": time.Since(start).String(),
	}).Info("redis del information")

	return errors.WithStack(err)
}

func (c *Client) HSet(ctx context.Context, key string, values ...interface{}) error {
	start := time.Now()
	err := c.Client.HSet(ctx, key, values...).Err()

	logx.WithContext(ctx).WithFields(logrus.Fields{
		"key":      key,
		"duration": time.Since(start).String(),
	}).Info("redis hset information")

	return errors.WithStack(err)
}

func (c *Client) HGet(ctx context.Context, key, field string) *Bind {
	start := time.Now()
	val, err := c.Client.HGet(ctx, key, field).Result()

	logx.WithContext(ctx).WithFields(logrus.Fields{
		"key":      key,
		"value":    val,
		"duration": time.Since(start).String(),
	}).Info("redis hget information")

	return &Bind{
		Val: val,
		Err: errors.WithStack(err),
	}
}

func (c *Client) HGetAll(ctx context.Context, key string) *redis.StringStringMapCmd {
	start := time.Now()
	cmd := c.Client.HGetAll(ctx, key)

	logx.WithContext(ctx).WithFields(logrus.Fields{
		"key":      key,
		"value":    cmd.Val(),
		"duration": time.Since(start).String(),
	}).Info("redis hgetall information")

	return cmd
}

func (c *Client) HDel(ctx context.Context, key string, fields ...string) error {
	start := time.Now()
	err := c.Client.HDel(ctx, key, fields...).Err()

	logx.WithContext(ctx).WithFields(logrus.Fields{
		"key":      key,
		"duration": time.Since(start).String(),
	}).Info("redis hdel information")

	return errors.WithStack(err)
}

func (c *Client) GetSet(ctx context.Context, key string, value interface{}) error {
	start := time.Now()
	err := c.Client.GetSet(ctx, key, value).Err()

	logx.WithContext(ctx).WithFields(logrus.Fields{
		"key":      key,
		"values":   value,
		"duration": time.Since(start).String(),
	}).Info("redis getset information")

	return errors.WithStack(err)
}

type Bind struct {
	Val string
	Err error
}

func (b *Bind) Bind(i interface{}) error {
	if b.Err != nil {
		return b.Err
	}

	err := json.Unmarshal([]byte(b.Val), i)
	return errors.WithStack(err)
}
