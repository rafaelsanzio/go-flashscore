package store

import (
	"context"
	"log"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/store/cursor"
	"github.com/rafaelsanzio/go-flashscore/pkg/store/mongo"
	"github.com/rafaelsanzio/go-flashscore/pkg/store/query"
)

// Simple CRUD interface
type Store interface {
	FindOne(ctx context.Context, collection string, filter query.Filter, v interface{}, opts ...query.FindOneOptions) errs.AppError
	Find(ctx context.Context, collection string, filter query.Filter, opts ...query.FindOptions) (cursor.Cursor, errs.AppError)
	InsertOne(ctx context.Context, collection string, data interface{}) (string, errs.AppError)
	UpdateOne(ctx context.Context, collection string, data interface{}) errs.AppError
	DeleteOne(ctx context.Context, collection string, id string) errs.AppError
}

var store Store

func init() {
	s, err := mongo.NewStore()
	if err != nil {
		log.Fatalf("error creating mongo store: %v", err)
	}

	store = s
}

func GetStore() Store {
	return store
}

func SetStore(s Store) {
	store = s
}
