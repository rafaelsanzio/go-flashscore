package cursor

import (
	"context"
)

type Cursor interface {
	Next(ctx context.Context) bool
	Decode(v interface{}) error
	Close(ctx context.Context) error
	Err() error
}
