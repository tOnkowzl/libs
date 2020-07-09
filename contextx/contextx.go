package contextx

import "context"

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
