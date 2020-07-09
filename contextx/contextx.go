package contextx

import (
	"context"

	"github.com/google/uuid"
)

type Key int

const (
	ID Key = iota
)

func GetID(ctx context.Context) string {
	if id, ok := ctx.Value(ID).(string); ok {
		return id
	}

	return ""
}

func WithID() context.Context {
	return context.WithValue(context.Background(), ID, uuid.New().String())
}
