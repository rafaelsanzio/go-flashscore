package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/rafaelsanzio/go-flashscore/pkg/config"
	"github.com/rafaelsanzio/go-flashscore/pkg/config/key"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
)

var jwtParseWithClaims = jwt.ParseWithClaims

func fakeJwtParseWithClaims(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error) {
	return nil, errs.ErrMarshalingJson
}

func restoreJwtParseWithClaims(replace func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error)) {
	jwtParseWithClaims = replace
}

var configValue = config.Value

func fakeConfigValue(k key.Key) (string, errs.AppError) {
	return "", errs.ErrMarshalingJson
}

func fakeConfigValueWithAnyKey(k key.Key) (string, errs.AppError) {
	return "any_key", nil
}

func restoreConfigValue(replace func(k key.Key) (string, errs.AppError)) {
	configValue = replace
}
