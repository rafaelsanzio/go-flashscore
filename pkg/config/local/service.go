package local

import (
	"github.com/rafaelsanzio/go-flashscore/pkg/config/key"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
)

type Service struct{}

func (s Service) Value(key key.Key) (string, errs.AppError) {
	return Values[key], nil
}
