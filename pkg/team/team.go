package team

import (
	"context"
	"time"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
)

type TeamRepo interface {
	Insert(ctx context.Context, t Team) errs.AppError
	Get(ctx context.Context, id string) (*Team, errs.AppError)
	List(ctx context.Context) ([]Team, errs.AppError)
	Update(ctx context.Context, t Team) (*Team, errs.AppError)
	Delete(ctx context.Context, id string) errs.AppError
}

type Team struct {
	ID        string `bson:"_id"`
	Name      string
	ShortCode string
	Country   string
	City      string
	Created   time.Time
}

func (t Team) GetID() string {
	return t.ID
}

func (t *Team) SetID(id string) {
	t.ID = id
}
