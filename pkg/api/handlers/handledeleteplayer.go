package handlers

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
)

func HandleDeletePlayer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		_ = errs.ErrGettingParam.Throwf(applog.Log, errs.ErrFmt, id)
		errs.HttpNotFound(w)
		return
	}

	player, err := repo.GetPlayerRepo().Get(ctx, id)
	if err != nil {
		errs.HttpInternalServerError(w)
		return
	}

	if player == nil {
		_ = errs.ErrPlayerIsNotFound.Throwf(applog.Log, errs.ErrFmt, id)
		errs.HttpNotFound(w)
		return
	}

	err = repo.GetPlayerRepo().Delete(ctx, id)
	if err != nil {
		errs.HttpInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
