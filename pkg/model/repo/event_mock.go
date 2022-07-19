package repo

import (
	"context"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/event"
)

type MockEventRepo struct {
	event.EventRepo
	InsertFunc              func(ctx context.Context, mt event.Event) errs.AppError
	ListEventsFromMatchFunc func(ctx context.Context, matchID string) ([]event.Event, errs.AppError)
}

func (m MockEventRepo) Insert(ctx context.Context, mt event.Event) errs.AppError {
	if m.InsertFunc != nil {
		return m.InsertFunc(ctx, mt)
	}
	return m.EventRepo.Insert(ctx, mt)
}

func (m MockEventRepo) ListEventsFromMatch(ctx context.Context, matchID string) ([]event.Event, errs.AppError) {
	if m.ListEventsFromMatchFunc != nil {
		return m.ListEventsFromMatchFunc(ctx, matchID)
	}
	return m.EventRepo.ListEventsFromMatch(ctx, matchID)
}
