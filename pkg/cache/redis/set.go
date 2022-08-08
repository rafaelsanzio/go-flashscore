package redis

import (
	"context"
	"time"
)

func (s *Store) Set(ctx context.Context, key string, value []byte) error {
	err := s.client.Set(ctx, key, value, 10*time.Second).Err()
	if err != nil {
		return err
	}

	return nil
}
