package hashmap

import (
	"context"
)

// Handler ...
type Handler func(key Key, value interface{}) bool

// WithInterruption ...
func WithInterruption(ctx context.Context, handler Handler) Handler {
	return func(key Key, value interface{}) bool {
		select {
		case <-ctx.Done():
			return false
		default:
			return handler(key, value)
		}
	}
}
