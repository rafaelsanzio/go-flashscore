package config

import (
	"github.com/rafaelsanzio/go-flashscore/pkg/config/key"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
)

type Service interface {
	Value(key.Key) (string, errs.AppError)
}
