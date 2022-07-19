package repo

import (
	"context"
	"time"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/event"
	"github.com/rafaelsanzio/go-flashscore/pkg/store"
	"github.com/rafaelsanzio/go-flashscore/pkg/store/query"
)

const (
	EventCollection = "event"
)

type eventRepo struct {
	store store.Store
}

var eventRepoSingleton event.EventRepo

func GetEventRepo() event.EventRepo {
	if eventRepoSingleton == nil {
		return getEventRepo()
	}
	return eventRepoSingleton
}

func getEventRepo() *eventRepo {
	s := store.GetStore()
	return &eventRepo{s}
}

func SetEventRepo(repo event.EventRepo) {
	eventRepoSingleton = repo
}

func (repo eventRepo) Insert(ctx context.Context, e event.Event) errs.AppError {
	e.Created = time.Now()
	_, err := repo.store.InsertOne(ctx, EventCollection, &e)
	return err
}

func (repo eventRepo) ListEventsFromMatch(ctx context.Context, matchID string) ([]event.Event, errs.AppError) {
	filter := query.Filter{
		"matchid": matchID,
	}

	opts := query.FindOptions{}
	mEvent := []event.Event{}
	events, err := repo.store.Find(ctx, EventCollection, filter, opts)
	if err != nil {
		return mEvent, errs.ErrMongoFind.Throwf(applog.Log, "for collection: %s, err: [%v]", EventCollection, err)
	}

	defer func() {
		_ = events.Close(ctx)
	}()

	for {
		if events.Err() != nil {
			return mEvent, err
		}

		if ok := events.Next(ctx); !ok {
			break
		}

		var e event.Event
		if err_ := events.Decode(&e); err_ != nil {
			return mEvent, err
		}

		mEvent = append(mEvent, e)
	}

	return mEvent, nil
}
