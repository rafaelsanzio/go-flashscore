package handlers

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
)

func HandleDeleteTeam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		_ = errs.ErrGettingParam.Throwf(applog.Log, errs.ErrFmt, id)
		errs.HttpNotFound(w)
		return
	}

	team, err := repo.GetTeamRepo().Get(ctx, id)
	if err != nil {
		errs.HttpInternalServerError(w)
		return
	}

	if team == nil {
		_ = errs.ErrTeamIsNotFound.Throwf(applog.Log, errs.ErrFmt, id)
		errs.HttpNotFound(w)
		return
	}

	err = repo.GetTeamRepo().Delete(ctx, id)
	if err != nil {
		errs.HttpInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
