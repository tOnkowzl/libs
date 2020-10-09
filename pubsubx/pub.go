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
}

type Pub struct {
	Client    *pubsub.Client
	Timeout   time.Duration
	ProjectID string
}

func NewPub(cfg PubConfig) *Pub {
	if cfg.Client == nil || cfg.Timeout == 0 || cfg.ProjectID == "" {
		logx.Panic("pubsub client, timeout, topicID, projectID required")
	}

	return &Pub{
		Client:    cfg.Client,
		Timeout:   cfg.Timeout,
		ProjectID: cfg.ProjectID,
	}
}

func (p *Pub) Publish(ctx context.Context, topicID string, i interface{}) error {
	b, err := json.Marshal(i)
	if err != nil {
		return errors.WithStack(err)
	}

	logx.WithContext(ctx).WithFields(logrus.Fields{
		"topicID":   topicID,
		"projectID": p.ProjectID,
		"value":     logx.LimitMSG(b),
	}).Info("pub information")

	pubCtx, cancel := context.WithTimeout(context.Background(), p.Timeout)
	defer cancel()

	topic := p.Client.Topic(topicID)

	if _, err = topic.Publish(pubCtx, &pubsub.Message{Data: b}).Get(pubCtx); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
