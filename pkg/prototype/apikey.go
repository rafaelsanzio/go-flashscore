package prototype

import (
	"github.com/rafaelsanzio/go-flashscore/pkg/apikey"
)

func PrototypeApiKey() apikey.ApiKey {
	return apikey.ApiKey{
		ID:    "1",
		Email: "rafael@gmail.com",
		Key:   "any_key",
	}
}
