package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/match"
	"github.com/rafaelsanzio/go-flashscore/pkg/model"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
)

func HandlePostMatch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		_ = errs.ErrGettingParam.Throwf(applog.Log, errs.ErrFmt, id)
		errs.HttpNotFound(w)
		return
	}

	tournament, err := repo.GetTournamentRepo().Get(ctx, id)
	if err != nil {
		errs.HttpInternalServerError(w)
		return
	}

	if tournament == nil {
		_ = errs.ErrTournamentIsNotFound.Throwf(applog.Log, errs.ErrFmt, id)
		errs.HttpNotFound(w)
		return
	}

	matchPayload, err := decodeMatchRequest(r)
	if err != nil {
		errs.HttpUnprocessableEntity(w, fmt.Sprintf("err: [%v]", err.Error()))
		return
	}

	match, err := convertAndValidatePayloadToMatch(ctx, matchPayload)
	if err != nil || match == nil {
		errs.HttpUnprocessableEntity(w, fmt.Sprintf("err: [%v]", err.Error()))
		return
	}

	match.Tournament = *tournament
	match.Status = model.MatchStatusNotStart

	err = repo.GetMatchRepo().Insert(ctx, *match)
	if err != nil {
		errs.HttpInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func decodeMatchRequest(r *http.Request) (MatchEntityPayload, errs.AppError) {
	payload := MatchEntityPayload{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		return payload, errs.ErrUnmarshalingJson.Throwf(applog.Log, errs.ErrFmt, err)
	}

	return payload, nil
}

func convertAndValidatePayloadToMatch(ctx context.Context, mt MatchEntityPayload) (*match.Match, errs.AppError) {
	homeTeam, err := repo.GetTeamRepo().Get(ctx, mt.HomeTeam)
	if err != nil || homeTeam == nil {
		return nil, errs.ErrTeamIsNotFound.Throwf(applog.Log, errs.ErrFmt, mt.HomeTeam)
	}

	awayTeam, err := repo.GetTeamRepo().Get(ctx, mt.AwayTeam)
	if err != nil || awayTeam == nil {
		return nil, errs.ErrTeamIsNotFound.Throwf(applog.Log, errs.ErrFmt, mt.AwayTeam)
	}

	if homeTeam.ID == awayTeam.ID {
		return nil, errs.ErrValidation.Throwf(applog.Log, "The same teams cannot do a match")
	}

	_, err_ := timeParse("2006-01-02", mt.DateOfMatch)
	if err_ != nil {
		return nil, errs.ErrParsingTime.Throwf(applog.Log, errs.ErrFmt, err_.Error())
	}

	_, err_ = timeParse("15:04", mt.TimeOfMatch)
	if err_ != nil {
		return nil, errs.ErrParsingTime.Throwf(applog.Log, errs.ErrFmt, err_.Error())
	}

	return &match.Match{
		HomeTeam:    *homeTeam,
		AwayTeam:    *awayTeam,
		DateOfMatch: mt.DateOfMatch,
		TimeOfMatch: mt.TimeOfMatch,
	}, nil
}
