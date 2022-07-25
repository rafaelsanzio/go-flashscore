package handlers

import (
	"encoding/json"
	"net/http"
	"net/mail"

	"github.com/rafaelsanzio/go-flashscore/pkg/apikey"
	"github.com/rafaelsanzio/go-flashscore/pkg/config"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
)

var rateLimitAllow = rateLimit.Allow

func fakeRateLimitAllow() bool {
	return false
}

func restoreRateLimitAllow(replace func() bool) {
	rateLimitAllow = replace
}

var jsonMarshal = json.Marshal

func fakeMarshal(v interface{}) ([]byte, error) {
	return []byte{}, errs.ErrMarshalingJson
}

func restoreMarshal(replace func(v interface{}) ([]byte, error)) {
	jsonMarshal = replace
}

var write = http.ResponseWriter.Write

func fakeWrite(http.ResponseWriter, []byte) (int, error) {
	return 0, errs.ErrResponseWriter
}

func restoreWrite(replace func(http.ResponseWriter, []byte) (int, error)) {
	write = replace
}

var convertPayloadToApiKeyFunc = convertPayloadToApiKey

func fakeConvertPayloadToApiKeyFunc(a APIKeyPayload) (apikey.ApiKey, errs.AppError) {
	return apikey.ApiKey{}, errs.ErrConvertingPayload
}

func restoreConvertPayloadToApiKeyFunc(replace func(a APIKeyPayload) (apikey.ApiKey, errs.AppError)) {
	convertPayloadToApiKeyFunc = replace
}

var configGenerateJWTToken = config.GenerateJWTToken

func fakeConfigGenerateJWTToken(email string) (string, errs.AppError) {
	return "", errs.ErrGenerateJWT
}

func restoreConfigGenerateJWTToken(replace func(email string) (string, errs.AppError)) {
	configGenerateJWTToken = replace
}

func validMailAddress(address string) bool {
	_, err := mail.ParseAddress(address)
	return err == nil
}
