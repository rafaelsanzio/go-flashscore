package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
	"github.com/rafaelsanzio/go-flashscore/pkg/tournament"
)

func HandleAddTeamsTournament(w http.ResponseWriter, r *http.Request) {
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
		_ = errs.ErrPlayerIsNotFound.Throwf(applog.Log, errs.ErrFmt, id)
		errs.HttpNotFound(w)
		return
	}

	tournamentPayload, err := decodeAddTeamsTournamentRequest(r)
	if err != nil {
		errs.HttpUnprocessableEntity(w, fmt.Sprintf("err: [%v]", err.Error()))
		return
	}

	if len(tournamentPayload.Teams) == 0 {
		errs.HttpUnprocessableEntity(w, "err: [Teams array cannot be null]")
		return
	}

	err = addTeamsToTournament(ctx, *tournament, tournamentPayload.Teams)
	if err != nil {
		errs.HttpInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func addTeamsToTournament(ctx context.Context, t tournament.Tournament, teamsID []string) errs.AppError {
	for _, teamID := range teamsID {
		team, err := repo.GetTeamRepo().Get(ctx, teamID)
		if err != nil {
			return errs.ErrTeamIsNotFound.Throwf(applog.Log, errs.ErrFmt, teamID)
		}
		if !t.FindTeam(teamID) {
			t.Teams = append(t.Teams, *team)
		}
	}

	_, err := repo.GetTournamentRepo().Update(ctx, t)
	if err != nil {
		return err
	}

	return nil
}

func decodeAddTeamsTournamentRequest(r *http.Request) (AddTeamsTournamentEntityPayload, errs.AppError) {
	payload := AddTeamsTournamentEntityPayload{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		return payload, errs.ErrUnmarshalingJson.Throwf(applog.Log, errs.ErrFmt, err)
	}

	return payload, nil
}
