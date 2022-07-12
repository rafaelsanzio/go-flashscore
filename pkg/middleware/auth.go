package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/config"
	"github.com/rafaelsanzio/go-flashscore/pkg/config/key"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
)

func IsAuthorized(ctx context.Context, token string) errs.AppError {
	if token == "" {
		return errs.ErrTokenIsEmpty.Throwf(applog.Log, errs.ErrFmt)
	}

	tokenString := strings.Split(token, " ")
	token = tokenString[1]

	secretApiKey, err_ := configValue(key.SecretApiKey)
	if err_ != nil {
		return err_
	}

	// Initialize a new instance of `Claims`
	claims := &config.Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretApiKey), nil
	})
	if err != nil {
		return errs.ErrTokenInvalid.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}

	if !tkn.Valid {
		return errs.ErrTokenInvalid.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}

	apikey, err_ := repo.GetApiKeyRepo().FindByEmail(ctx, claims.Email)
	if err_ != nil {
		return err_
	}
	if apikey == nil {
		return errs.ErrTokenWithInvalidEmail.Throwf(applog.Log, errs.ErrFmt, claims.Email)
	}

	validKey := apikey.ValidUntil.After(time.Now())
	if !validKey {
		return errs.ErrTokenInvalid.Throwf(applog.Log, errs.ErrFmt, claims.Email)
	}

	return nil
}
