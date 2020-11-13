package saramax

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tOnkowzl/libs/logx"
)

type Config struct {
	*sarama.Config
	Brokers []string
}

type Produce struct {
	SyncProducer sarama.SyncProducer
}

func NewProduce(cfg *Config) *Produce {
	syncProducer, err := sarama.NewSyncProducer(cfg.Brokers, cfg.Config)
	if err != nil {
		logx.Panic(err)
	}

	return &Produce{
		SyncProducer: syncProducer,
	}
}

func (p *Produce) Produce(ctx context.Context, topic string, i interface{}) error {
	b, err := json.Marshal(i)
	if err != nil {
		return errors.WithStack(err)
	}

	start := time.Now()
	partition, offset, err := p.SyncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(b),
	})

	logx.WithContext(ctx).WithFields(logrus.Fields{
		"value":     logx.LimitMSGByte(b),
		"topic":     topic,
		"partition": partition,
		"offset":    offset,
		"duration":  time.Since(start).String(),
	}).Info("produce information")

	return errors.WithStack(err)
}
