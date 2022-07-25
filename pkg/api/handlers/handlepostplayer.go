package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
	"github.com/rafaelsanzio/go-flashscore/pkg/player"
)

func HandlePostPlayer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	err = repo.GetPlayerRepo().Insert(ctx, *player)
	if err != nil {
		errs.HttpInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func decodePlayerRequest(r *http.Request) (PlayerEntityPayload, errs.AppError) {
	payload := PlayerEntityPayload{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		return payload, errs.ErrUnmarshalingJson.Throwf(applog.Log, errs.ErrFmt, err)
	}

	return payload, nil
}

func convertPayloadToPlayer(ctx context.Context, p PlayerEntityPayload) (*player.Player, errs.AppError) {
	team, err := repo.GetTeamRepo().Get(ctx, p.Team)
	if err != nil {
		return nil, errs.ErrConvertingPayload.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}
	if team == nil {
		return nil, errs.ErrTeamIsNotFound.Throwf(applog.Log, errs.ErrFmt, p.Team)
	}

	result := player.Player{
		Name:         p.Name,
		Team:         *team,
		Country:      p.Country,
		BirthdayDate: p.BirthdayDate,
	}

	return &result, nil
}
