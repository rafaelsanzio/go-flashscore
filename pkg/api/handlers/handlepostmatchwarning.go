package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/event"
	"github.com/rafaelsanzio/go-flashscore/pkg/kafka"
	"github.com/rafaelsanzio/go-flashscore/pkg/model"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
	"github.com/rafaelsanzio/go-flashscore/pkg/player"
	"github.com/rafaelsanzio/go-flashscore/pkg/team"
)

func HandlePostMatchWarning(w http.ResponseWriter, r *http.Request) {
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

	matchID := vars["match_id"]

	if matchID == "" {
		_ = errs.ErrGettingParam.Throwf(applog.Log, errs.ErrFmt, id)
		errs.HttpNotFound(w)
		return
	}

	match, err := repo.GetMatchRepo().FindMatchForTournament(ctx, matchID, tournament.ID)
	if err != nil {
		errs.HttpInternalServerError(w)
		return
	}

	if match == nil {
		_ = errs.ErrMatchIsNotFound.Throwf(applog.Log, errs.ErrFmt, id)
		errs.HttpNotFound(w)
		return
	}

	if match.Status != model.MatchStatusInProgress {
		errs.HttpUnprocessableEntity(w, fmt.Sprintf("err: [%v]", "A warning cannot happen when the game is not in progress"))
		return
	}

	matchWarningPayload, err := decodeMatchWarningRequest(r)
	if err != nil {
		errs.HttpUnprocessableEntity(w, fmt.Sprintf("err: [%v]", err.Error()))
		return
	}

	teamInMatch := match.FindTeamInMatch(matchWarningPayload.Team)
	if !teamInMatch {
		errs.HttpUnprocessableEntity(w, fmt.Sprintf("this team is not in this match: [%v]", matchWarningPayload.Team))
		return
	}

	teamWarn, playerWarn, err := convertAndValidatePayloadToMatchWarning(ctx, matchWarningPayload)
	if err != nil {
		errs.HttpUnprocessableEntity(w, fmt.Sprintf("err: [%v]", err.Error()))
		return
	}

	minute := matchWarningPayload.Minute
	minuteAsString := strconv.Itoa(minute)

	value := struct {
		Team          team.Team
		Player        player.Player
		Warning       model.Warnings
		WarningMinute int
		Created       time.Time
	}{
		Team:          *teamWarn,
		Player:        *playerWarn,
		Warning:       matchWarningPayload.Warning,
		WarningMinute: minute,
		Created:       time.Now(),
	}

	event := event.Event{
		TournamentID: tournament.ID,
		MatchID:      match.ID,
		Type:         model.EventWarning,
		Value:        value,
	}

	err = repo.GetEventRepo().Insert(ctx, event)
	if err != nil {
		errs.HttpInternalServerError(w)
		return
	}

	data := map[string]string{
		"matchEventType": string(model.EventWarning),
		"tournamentID":   tournament.ID,
		"matchID":        match.ID,
		"teamID":         teamWarn.ID,
		"playerID":       playerWarn.ID,
		"warning":        string(matchWarningPayload.Warning),
		"warningMinute":  minuteAsString,
	}

	go kafka.Notify(ctx, data, model.ActionGameEvents, "Game Events - Warning Player", model.KafkaTopicMatchEvents)

	w.WriteHeader(http.StatusCreated)
}

func decodeMatchWarningRequest(r *http.Request) (MatchWarningPayload, errs.AppError) {
	payload := MatchWarningPayload{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		return payload, errs.ErrUnmarshalingJson.Throwf(applog.Log, errs.ErrFmt, err)
	}

	return payload, nil
}

func convertAndValidatePayloadToMatchWarning(ctx context.Context, mt MatchWarningPayload) (*team.Team, *player.Player, errs.AppError) {
	team, err := repo.GetTeamRepo().Get(ctx, mt.Team)
	if err != nil || team == nil {
		return nil, nil, errs.ErrTeamIsNotFound.Throwf(applog.Log, errs.ErrFmt, mt.Team)
	}

	player, err := repo.GetPlayerRepo().GetTeamPlayer(ctx, mt.Player, team.ID)
	if err != nil || player == nil {
		return nil, nil, errs.ErrPlayerIsNotFoundInThisTeam.Throwf(applog.Log, errs.ErrFmt, team.ID, mt.Player)
	}

	return team, player, nil
}
