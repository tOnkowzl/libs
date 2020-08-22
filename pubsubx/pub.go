package pubsubx

import (
	"context"
	"encoding/json"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tOnkowzl/libs/logx"
)

type PubConfig struct {
	Client    *pubsub.Client
	Timeout   time.Duration
	ProjectID string
	TopicID   string
}

type Pub struct {
	Topic     *pubsub.Topic
	Timeout   time.Duration
	TopicID   string
	ProjectID string
}

func NewPub(cfg PubConfig) *Pub {
	if cfg.Client == nil || cfg.Timeout == 0 || cfg.TopicID == "" || cfg.ProjectID == "" {
		logx.Panic("pubsub client, timeout, topicID, projectID required")
	}

	return &Pub{
		Topic:     cfg.Client.Topic(cfg.TopicID),
		Timeout:   cfg.Timeout,
		TopicID:   cfg.TopicID,
		ProjectID: cfg.ProjectID,
	}
}

func (p *Pub) Publish(ctx context.Context, i interface{}) error {
	b, err := json.Marshal(i)
	if err != nil {
		return errors.WithStack(err)
	}

	logx.WithContext(ctx).WithFields(logrus.Fields{
		"topicID":   p.TopicID,
		"projectID": p.ProjectID,
		"value":     logx.LimitMSG(b),
	}).Info("pub information")

	pubCtx, cancel := context.WithTimeout(context.Background(), p.Timeout)
	defer cancel()

	res := p.Topic.Publish(pubCtx, &pubsub.Message{Data: b})
	if _, err = res.Get(pubCtx); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
