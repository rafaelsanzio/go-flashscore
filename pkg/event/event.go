package event

import (
	"context"
	"time"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/model"
)

type EventRepo interface {
	Insert(ctx context.Context, e Event) errs.AppError
	ListEventsFromMatch(ctx context.Context, matchID string) ([]Event, errs.AppError)
}

type Event struct {
	ID           string `bson:"_id"`
	TournamentID string
	MatchID      string
	Type         model.EventsMatchType
	Value        interface{}
	Created      time.Time
}

func (e Event) GetID() string {
	return e.ID
}

func (e *Event) SetID(id string) {
	e.ID = id
}
