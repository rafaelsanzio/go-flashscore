package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/kafka"
	"github.com/rafaelsanzio/go-flashscore/pkg/model"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
	"github.com/rafaelsanzio/go-flashscore/pkg/transfer"
)

func HandlePostTransfer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	transferPayload, err := decodeTransferRequest(r)
	if err != nil {
		errs.HttpUnprocessableEntity(w, fmt.Sprintf("err: [%v]", err.Error()))
		return
	}

	transfer, err := convertAndValidatePayloadToTransferFunc(ctx, transferPayload)
	if err != nil {
		errs.HttpUnprocessableEntity(w, fmt.Sprintf("err: [%v]", err.Error()))
		return
	}

	err = repo.GetTransferRepo().Insert(ctx, *transfer)
	if err != nil {
		errs.HttpInternalServerError(w)
		return
	}

	data := map[string]string{
		"playerID":      transfer.Player.ID,
		"teamDestinyID": transfer.TeamDestiny.ID,
	}

	go kafka.Notify(ctx, data, model.ActionUpdateTeamPlayer, "Transfer Player", model.KafkaTopicTransfer)

	w.WriteHeader(http.StatusCreated)
}

func decodeTransferRequest(r *http.Request) (TransferEntityPayload, errs.AppError) {
	payload := TransferEntityPayload{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		return payload, errs.ErrUnmarshalingJson.Throwf(applog.Log, errs.ErrFmt, err)
	}

	return payload, nil
}

func convertAndValidatePayloadToTransfer(ctx context.Context, t TransferEntityPayload) (*transfer.Transfer, errs.AppError) {
	teamDestiny, err := repo.GetTeamRepo().Get(ctx, t.TeamDestiny)
	if err != nil {
		return nil, errs.ErrConvertingPayload.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}
	if teamDestiny == nil {
		return nil, errs.ErrTeamIsNotFound.Throwf(applog.Log, errs.ErrFmt, t.TeamDestiny)
	}

	player, err := repo.GetPlayerRepo().Get(ctx, t.Player)
	if err != nil {
		return nil, errs.ErrConvertingPayload.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}
	if player == nil {
		return nil, errs.ErrPlayerIsNotFound.Throwf(applog.Log, errs.ErrFmt, t.Player)
	}

	_, err_ := timeParse("2006-01-02", t.DateOfTransfer)
	if err_ != nil {
		return nil, errs.ErrParsingTime.Throwf(applog.Log, errs.ErrFmt, err_.Error())
	}

	if player.Team.ID == teamDestiny.ID {
		return nil, errs.ErrValidation.Throwf(applog.Log, "Player cannot be transfer from: %s to: %s  is the same team", player.Team.ID, teamDestiny.ID)
	}

	result := transfer.Transfer{
		Player:         *player,
		TeamDestiny:    *teamDestiny,
		DateOfTransfer: t.DateOfTransfer,
		Amount:         t.Amount,
	}

	return &result, nil
}
