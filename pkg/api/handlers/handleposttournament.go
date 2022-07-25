package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
	"github.com/rafaelsanzio/go-flashscore/pkg/team"
	"github.com/rafaelsanzio/go-flashscore/pkg/tournament"
)

func HandlePostTournament(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tournamentPayload, err := decodeTournamentRequest(r)
	if err != nil {
		errs.HttpUnprocessableEntity(w, fmt.Sprintf("err: [%v]", err.Error()))
		return
	}

	tournament, err := convertPayloadToTournamentFunc(ctx, tournamentPayload)
	if err != nil {
		errs.HttpUnprocessableEntity(w, fmt.Sprintf("err: [%v]", err.Error()))
		return
	}

	err = repo.GetTournamentRepo().Insert(ctx, *tournament)
	if err != nil {
		errs.HttpInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func decodeTournamentRequest(r *http.Request) (TournamentEntityPayload, errs.AppError) {
	payload := TournamentEntityPayload{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		return payload, errs.ErrUnmarshalingJson.Throwf(applog.Log, errs.ErrFmt, err)
	}

	return payload, nil
}

func convertPayloadToTournament(ctx context.Context, t TournamentEntityPayload) (*tournament.Tournament, errs.AppError) {
	var teams []team.Team
	if len(t.Teams) > 0 {
		for _, teamID := range t.Teams {
			team, err := repo.GetTeamRepo().Get(ctx, teamID)
			if err != nil || team == nil {
				return nil, errs.ErrTeamIsNotFound.Throwf(applog.Log, errs.ErrFmt, teamID)
			}
			teams = append(teams, *team)
		}
	}

	result := tournament.Tournament{
		Name:  t.Name,
		Teams: teams,
	}

	return &result, nil
}
