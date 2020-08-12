package pubsubx

import (
	"context"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/tOnkowzl/libs/logx"
	"google.golang.org/api/option"
)

type ClientConfig struct {
	Pool        int
	Timeout     time.Duration
	ProjectID   string
	Credentials string
}

func NewClient(cfg *ClientConfig) *pubsub.Client {
	if cfg.Pool == 0 || cfg.Timeout == 0 || cfg.ProjectID == "" || cfg.Credentials == "" {
		logx.Panic("pubsub pool, timeout, projectID, credentails required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	client, err := pubsub.NewClient(ctx,
		cfg.ProjectID,
		option.WithCredentialsJSON([]byte(cfg.Credentials)),
		option.WithGRPCConnectionPool(cfg.Pool),
	)
	if err != nil {
		logx.Panic(err)
	}

	return client
}
