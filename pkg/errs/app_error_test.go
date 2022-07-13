package errs

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
)

func TestEnsureUnique(t *testing.T) {
	duplicateCode := "ABC123"
	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, r, fmt.Sprintf("AppError Code %v already exists", duplicateCode))
		} else {
			assert.NotNil(t, r)
		}
	}()

	_ = _new(duplicateCode, "err1")
	_ = _new(duplicateCode, "err2")
}

func TestEnsureFormatWithBadCode(t *testing.T) {
	testCases := []struct {
		Code string
	}{
		{Code: "ABC1234"},
		{Code: "ABCD123"},
		{Code: "abc123"},
		{Code: ""},
	}

	for _, tc := range testCases {
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, r, fmt.Sprintf("AppError Code must have format ABC123 but is %s", tc.Code))
			} else {
				assert.NotNil(t, r)
			}
		}()

		_ = _new(tc.Code, "err")
	}
}

func TestAppErrorAsInterface(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			assert.EqualError(t, err, "CMN000: error loading timezone data")
		} else {
			assert.NotNil(t, r)
		}
	}()

	panic(ErrLoadingTimeZone)
}

func TestAppErrorCode(t *testing.T) {
	log := applog.Log
	err := ErrLoadingTimeZone.Throw(log)
	te := err.(*appError)
	var correct bool
	switch te.Code() {
	case ErrLoadingTimeZone.Code():
		correct = true
	default:
		correct = false
	}
	assert.True(t, correct)
}

func TestAppErrorError(t *testing.T) {
	newError := _new("TST000", "newError")
	expected := "TST000: newError"
	actual := newError.Error()
	assert.Equal(t, expected, actual)
}

func TestAppErrorThrow(t *testing.T) {
	newError := _new("TST001", "newError")

	// test nil ctx
	nilCtxError := newError.Throwf(nil, "").(*appError)
	// test new error
	assert.False(t, newError == nilCtxError)

	log := applog.Log

	// test new msg/msg order
	newerError := newError.Throw(log).(*appError)

	expected := "TST001: newError"
	actual := newerError.Error()
	assert.Equal(t, expected, actual)
	assert.False(t, newError == newerError)
}

func TestAppErrorThrowf(t *testing.T) {
	newError := _new("TST002", "newError")

	// test nil ctx
	nilCtxError := newError.Throwf(nil, "")
	// test new error
	assert.False(t, newError == nilCtxError)

	log := applog.Log

	// test new msg/msg order
	newerError := newError.Throwf(log, "is new")

	expected := "TST002: newError, is new"
	actual := newerError.Error()
	assert.Equal(t, expected, actual)
	assert.NotEqual(t, newError, newerError)

	externalError := errors.New("external error")
	appError := _new("TST003", "app error")
	newAppError := appError.Throwf(log, "err: %v", externalError)
	expected2 := "TST003: app error, err: external error"
	actual2 := newAppError.Error()
	assert.Equal(t, expected2, actual2)
	assert.NotEqual(t, appError, newAppError)
}

func TestAnnotatef(t *testing.T) {
	log := applog.Log
	err1 := _new("TST004", "err1")
	err2 := err1.Throw(log)
	err2 = err2.Annotatef(log, "extra %v", "data")
	err2 = err2.Annotatef(log, "%s %s", "more", "data")
	result := err2.Annotatef(log, "and more data")

	actual := err2.Error()
	assert.Equal(t, actual, result.Error())

	expected := "TST004: err1 - extra data, more data, and more data"
	assert.Equal(t, expected, actual)
}

func TestIs(t *testing.T) {
	log := applog.Log
	err := ErrLoadingTimeZone.Throw(log)
	res := err.Is(ErrLoadingTimeZone)
	assert.True(t, res)
}
