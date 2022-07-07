package noop

import (
	"context"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/store/cursor"
	"github.com/rafaelsanzio/go-flashscore/pkg/store/query"
)

type Store struct{}

func (s Store) FindOne(_ context.Context, _ string, _ query.Filter, _ interface{}, _ ...query.FindOneOptions) errs.AppError {
	return nil
}

func (s Store) Find(_ context.Context, _ string, _ query.Filter, _ ...query.FindOptions) (cursor.Cursor, errs.AppError) {
	return Cursor{}, nil
}

func (s Store) InsertOne(_ context.Context, _ string, _ interface{}) (string, errs.AppError) {
	return "", nil
}

func (s Store) UpdateOne(_ context.Context, _ string, _ interface{}) errs.AppError {
	return nil
}

func (s Store) DeleteOne(_ context.Context, _ string, _ string) errs.AppError {
	return nil
}

type Cursor struct{}

func (c Cursor) Next(_ context.Context) bool {
	return false
}

func (c Cursor) Decode(_ interface{}) error {
	return nil
}

func (c Cursor) Close(_ context.Context) error {
	return nil
}

func (c Cursor) Err() error {
	return nil
}
