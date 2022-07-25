package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rafaelsanzio/go-flashscore/pkg/apikey"
	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
)

func HandleCreateAPIKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	APIKeyPayload, err := decodeAPIKeyRequest(r)
	if err != nil {
		errs.HttpUnprocessableEntity(w, fmt.Sprintf("err: [%v]", err.Error()))
		return
	}

	isMailValid := validMailAddress(APIKeyPayload.Email)
	if !isMailValid {
		err = errs.ErrInvalidEmail.Throwf(applog.Log, errs.ErrFmt, APIKeyPayload.Email)
		errs.HttpUnprocessableEntity(w, fmt.Sprintf("err: [%v]", err.Error()))
		return
	}

	ap, err := repo.GetApiKeyRepo().FindByEmail(ctx, APIKeyPayload.Email)
	if err != nil {
		errs.HttpInternalServerError(w)
		return
	}

	tokenString, err_ := configGenerateJWTToken(APIKeyPayload.Email)
	if err_ != nil {
		errs.HttpUnprocessableEntity(w, fmt.Sprintf("err: [%v]", err_.Error()))
		return
	}

	APIKeyPayload.Key = tokenString

	apikey, err := convertPayloadToApiKeyFunc(APIKeyPayload)
	if err != nil {
		errs.HttpUnprocessableEntity(w, fmt.Sprintf("err: [%v]", err.Error()))
		return
	}

	if ap != nil {
		ap.Key = apikey.Key
		ap.ValidUntil = apikey.ValidUntil
		ap, err = repo.GetApiKeyRepo().Update(ctx, *ap)
		if err != nil {
			errs.HttpInternalServerError(w)
			return
		}
		apikey = *ap
	} else {
		err = repo.GetApiKeyRepo().Insert(ctx, apikey)
		if err != nil {
			errs.HttpInternalServerError(w)
			return
		}
	}

	data, err__ := jsonMarshal(apikey)
	if err__ != nil {
		_ = errs.ErrMarshalingJson.Throwf(applog.Log, errs.ErrFmt, err_)
		errs.HttpInternalServerError(w)
		return
	}

	_, err__ = write(w, data)

	if err__ != nil {
		_ = errs.ErrResponseWriter.Throwf(applog.Log, errs.ErrFmt, err_)
		errs.HttpInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func decodeAPIKeyRequest(r *http.Request) (APIKeyPayload, errs.AppError) {
	payload := APIKeyPayload{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		return payload, errs.ErrUnmarshalingJson.Throwf(applog.Log, errs.ErrFmt, err)
	}

	return payload, nil
}

func convertPayloadToApiKey(a APIKeyPayload) (apikey.ApiKey, errs.AppError) {
	result := apikey.ApiKey{
		Email:      a.Email,
		Key:        a.Key,
		ValidUntil: time.Now().AddDate(0, 0, 10).Truncate(time.Second),
	}

	return result, nil
}
