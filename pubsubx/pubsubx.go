package pub

import (
	"context"
	"encoding/json"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tOnkowzl/libs/logx"
	"google.golang.org/api/option"
)

type Config struct {
	Timeout     time.Duration
	ProjectID   string
	TopicID     string
	Credentials string
}

type Pub struct {
	Topic     *pubsub.Topic
	Timeout   time.Duration
	TopicID   string
	ProjectID string
}

func NewPub(cfg Config) *Pub {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	client, err := pubsub.NewClient(ctx,
		cfg.ProjectID,
		option.WithCredentialsJSON([]byte(cfg.Credentials)))
	if err != nil {
		logx.Panic(err)
	}

	return &Pub{
		Topic:     client.Topic(cfg.TopicID),
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
		"value":     string(b),
	}).Info("publish information")

	pubCtx, cancel := context.WithTimeout(context.Background(), p.Timeout)
	defer cancel()

	res := p.Topic.Publish(pubCtx, &pubsub.Message{Data: b})
	if _, err = res.Get(pubCtx); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
