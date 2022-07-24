package handlers

import (
	"net/http"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
)

func HandleListPlayer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	player, err := repo.GetPlayerRepo().List(ctx)
	if err != nil {
		errs.HttpInternalServerError(w)
		return
	}

	data, err_ := jsonMarshal(player)
	if err_ != nil {
		_ = errs.ErrMarshalingJson.Throwf(applog.Log, errs.ErrFmt, err_)
		errs.HttpInternalServerError(w)
		return
	}

	_, err_ = write(w, data)
	if err_ != nil {
		_ = errs.ErrResponseWriter.Throwf(applog.Log, errs.ErrFmt, err_)
		errs.HttpInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
