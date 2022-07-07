package cursor

import (
	"context"
	"reflect"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
)

// ValueCursor is a simple cursor implementation for iterating through a collection of values, useful for mocking query responses
type ValueCursor struct {
	i      int
	values []interface{}
}

func NewValueCursor(values ...interface{}) *ValueCursor {
	return &ValueCursor{values: values, i: -1}
}

func (c *ValueCursor) Next(_ context.Context) bool {
	c.i++
	return c.i <= len(c.values)-1
}

func (c *ValueCursor) Decode(v interface{}) error {
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return errs.ErrDecodeCursor.Throw(applog.Log)
	}

	reflect.ValueOf(v).Elem().Set(reflect.ValueOf(c.values[c.i]))
	return nil
}

func (c ValueCursor) Close(_ context.Context) error {
	return nil
}

func (c ValueCursor) Err() error {
	return nil
}
