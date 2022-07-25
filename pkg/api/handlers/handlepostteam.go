package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
	"github.com/rafaelsanzio/go-flashscore/pkg/team"
)

func HandlePostTeam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	teamPayload, err := decodeTeamRequest(r)
	if err != nil {
		errs.HttpUnprocessableEntity(w, fmt.Sprintf("err: [%v]", err.Error()))
		return
	}

	team, err := convertPayloadToTeamFunc(teamPayload)
	if err != nil {
		errs.HttpUnprocessableEntity(w, fmt.Sprintf("err: [%v]", err.Error()))
		return
	}

	err = repo.GetTeamRepo().Insert(ctx, team)
	if err != nil {
		errs.HttpInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func decodeTeamRequest(r *http.Request) (TeamEntityPayload, errs.AppError) {
	payload := TeamEntityPayload{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		return payload, errs.ErrUnmarshalingJson.Throwf(applog.Log, errs.ErrFmt, err)
	}

	return payload, nil
}

func convertPayloadToTeam(t TeamEntityPayload) (team.Team, errs.AppError) {
	result := team.Team{
		Name:      t.Name,
		ShortCode: t.ShortCode,
		Country:   t.Country,
		City:      t.City,
	}

	return result, nil
}
