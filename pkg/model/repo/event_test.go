package repo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/event"
	"github.com/rafaelsanzio/go-flashscore/pkg/prototype"
)

func TestEventRepoInsert(t *testing.T) {
	ctx := context.Background()

	SetEventRepo(MockEventRepo{
		InsertFunc: func(ctx context.Context, mt event.Event) errs.AppError {
			return nil
		},
	})
	defer SetEventRepo(nil)

	newEvent := prototype.PrototypeEvent()

	err := GetEventRepo().Insert(ctx, newEvent)
	assert.NoError(t, err)
}

func TestEventRepoListEventsFromMatch(t *testing.T) {
	ctx := context.Background()

	SetEventRepo(MockEventRepo{
		ListEventsFromMatchFunc: func(ctx context.Context, matchID string) ([]event.Event, errs.AppError) {
			eventMock := prototype.PrototypeEvent()
			eventMock2 := prototype.PrototypeEvent()

			return []event.Event{eventMock, eventMock2}, nil
		},
	})
	defer SetEventRepo(nil)

	events, err := GetEventRepo().ListEventsFromMatch(ctx, "match-id")
	assert.NoError(t, err)

	assert.Equal(t, 2, len(events))
}
