package handlers

import (
	"strconv"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
)

var strconvAtoi = strconv.Atoi

func fakeStrconvAtoi(s string) (int, error) {
	return 0, errs.ErrParsingAtoi
}

func restoreStrconvAtoi(replace func(s string) (int, error)) {
	strconvAtoi = replace
}
