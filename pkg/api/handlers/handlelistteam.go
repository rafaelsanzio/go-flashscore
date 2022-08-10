package handlers

import (
	"fmt"
	"net/http"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/cache"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
)

func HandleListTeam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	team, err := repo.GetTeamRepo().List(ctx)
	if err != nil {
		errs.HttpInternalServerError(w)
		return
	}

	data, err_ := jsonMarshal(team)
	if err_ != nil {
		_ = errs.ErrMarshalingJson.Throwf(applog.Log, errs.ErrFmt, err_)
		errs.HttpInternalServerError(w)
		return
	}

	cacheKey := fmt.Sprintf("%s%s", r.Method, r.URL)
	_err := cache.GetStore().Set(ctx, cacheKey, data)
	if _err != nil {
		applog.Log.Warnf("Cache could not be set for this key: %s with error: %s", cacheKey, _err.Error())
	}

	_, err_ = write(w, data)
	if err_ != nil {
		_ = errs.ErrResponseWriter.Throwf(applog.Log, errs.ErrFmt, err_)
		errs.HttpInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
