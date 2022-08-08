package redis

import (
	"context"
	"encoding/json"
)

func (s *Store) Get(ctx context.Context, key string) (interface{}, error) {
	var thing interface{}
	thingAsBytes, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	if thingAsBytes != nil {
		err = json.Unmarshal(thingAsBytes, &thing)
		if err != nil {
			return nil, err
		}
	}

	return thing, nil
}
