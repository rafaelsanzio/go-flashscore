package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
)

func HandleUpdatePlayer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		_ = errs.ErrGettingParam.Throwf(applog.Log, errs.ErrFmt, id)
		errs.HttpNotFound(w)
		return
	}

	playerPayload, err := decodePlayerRequest(r)
	if err != nil {
		errs.HttpUnprocessableEntity(w, fmt.Sprintf("err: [%v]", err.Error()))
		return
	}

	player, err := convertPayloadToPlayerFunc(ctx, playerPayload)
	if err != nil {
		errs.HttpUnprocessableEntity(w, fmt.Sprintf("err: [%v]", err.Error()))
		return
	}

	player.ID = id
	_, err = repo.GetPlayerRepo().Update(ctx, *player)
	if err != nil {
		errs.HttpInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
