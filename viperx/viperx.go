package viperx

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tOnkowzl/libs/logx"
)

func log(ctx context.Context, key string, value interface{}) {
	logx.WithContext(ctx).WithFields(logrus.Fields{
		"key":   key,
		"value": value,
	}).Info("viper information")
}

func GetString(ctx context.Context, key string) string {
	value := viper.GetString(key)
	log(ctx, key, value)
	return value
}

func GetBool(ctx context.Context, key string) bool {
	value := viper.GetBool(key)
	log(ctx, key, value)
	return value
}

func GetInt(ctx context.Context, key string) int {
	value := viper.GetInt(key)
	log(ctx, key, value)
	return value
}

func GetInt32(ctx context.Context, key string) int32 {
	value := viper.GetInt32(key)
	log(ctx, key, value)
	return value
}

func GetInt64(ctx context.Context, key string) int64 {
	value := viper.GetInt64(key)
	log(ctx, key, value)
	return value
}

func GetUint(ctx context.Context, key string) uint {
	value := viper.GetUint(key)
	log(ctx, key, value)
	return value
}

func GetUint32(ctx context.Context, key string) uint32 {
	value := viper.GetUint32(key)
	log(ctx, key, value)
	return value
}

func GetUint64(ctx context.Context, key string) uint64 {
	value := viper.GetUint64(key)
	log(ctx, key, value)
	return value
}

func GetFloat64(ctx context.Context, key string) float64 {
	value := viper.GetFloat64(key)
	log(ctx, key, value)
	return value
}

func GetTime(ctx context.Context, key string) time.Time {
	value := viper.GetTime(key)
	log(ctx, key, value)
	return value
}

func GetDuration(ctx context.Context, key string) time.Duration {
	value := viper.GetDuration(key)
	log(ctx, key, value)
	return value
}
