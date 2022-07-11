package apikey

import (
	"context"
	"time"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
)

type ApiKeyRepo interface {
	Insert(ctx context.Context, a ApiKey) errs.AppError
	Get(ctx context.Context, id string) (*ApiKey, errs.AppError)
	FindByEmail(ctx context.Context, email string) (*ApiKey, errs.AppError)
	Update(ctx context.Context, a ApiKey) (*ApiKey, errs.AppError)
	Delete(ctx context.Context, id string) errs.AppError
}

type ApiKey struct {
	ID         string `bson:"_id"`
	Email      string
	Key        string
	Created    time.Time
	ValidUntil time.Time
}

func (a ApiKey) GetID() string {
	return a.ID
}

func (a *ApiKey) SetID(id string) {
	a.ID = id
}
