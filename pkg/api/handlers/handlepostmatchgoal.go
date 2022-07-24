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

func HandlePostMatchGoal(w http.ResponseWriter, r *http.Request) {
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
		errs.HttpUnprocessableEntity(w, fmt.Sprintf("err: [%v]", "Match cannot have a goal because it is not in progress"))
		return
	}

	matchGoalPayload, err := decodeMatchGoalRequest(r)
	if err != nil {
		errs.HttpUnprocessableEntity(w, fmt.Sprintf("err: [%v]", err.Error()))
		return
	}

	teamInMatch := match.FindTeamInMatch(matchGoalPayload.TeamScore)
	if !teamInMatch {
		errs.HttpUnprocessableEntity(w, fmt.Sprintf("this team is not in this match: [%v]", matchGoalPayload.TeamScore))
		return
	}

	teamScore, playerScore, goalMinute, err := convertAndValidatePayloadToMatchGoal(ctx, matchGoalPayload)
	if err != nil {
		errs.HttpUnprocessableEntity(w, fmt.Sprintf("err: [%v]", err.Error()))
		return
	}

	value := struct {
		TeamScore  team.Team
		Player     player.Player
		GoalMinute int
		Created    time.Time
	}{
		TeamScore:  *teamScore,
		Player:     *playerScore,
		GoalMinute: goalMinute,
		Created:    time.Now(),
	}

	event := event.Event{
		TournamentID: tournament.ID,
		MatchID:      match.ID,
		Type:         model.EventGoal,
		Value:        value,
	}

	err = repo.GetEventRepo().Insert(ctx, event)
	if err != nil {
		errs.HttpInternalServerError(w)
		return
	}

	goalMinuteAsString := strconv.Itoa(int(goalMinute))

	data := map[string]string{
		"matchEventType": string(model.EventGoal),
		"tournamentID":   tournament.ID,
		"matchID":        match.ID,
		"teamScore":      teamScore.ID,
		"player":         playerScore.ID,
		"goalMinute":     goalMinuteAsString,
	}

	go kafka.Notify(ctx, data, model.ActionGameEvents, "Game Events - Goal", model.KafkaTopicMatchEvents)

	w.WriteHeader(http.StatusCreated)
}

func decodeMatchGoalRequest(r *http.Request) (MatchGoalEntityPayload, errs.AppError) {
	payload := MatchGoalEntityPayload{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		return payload, errs.ErrUnmarshalingJson.Throwf(applog.Log, errs.ErrFmt, err)
	}

	return payload, nil
}

func convertAndValidatePayloadToMatchGoal(ctx context.Context, mt MatchGoalEntityPayload) (*team.Team, *player.Player, int, errs.AppError) {
	teamScore, err := repo.GetTeamRepo().Get(ctx, mt.TeamScore)
	if err != nil || teamScore == nil {
		return nil, nil, 0, errs.ErrTeamIsNotFound.Throwf(applog.Log, errs.ErrFmt, mt.TeamScore)
	}

	player, err := repo.GetPlayerRepo().GetTeamPlayer(ctx, mt.Player, teamScore.ID)
	if err != nil || player == nil {
		return nil, nil, 0, errs.ErrPlayerIsNotFoundInThisTeam.Throwf(applog.Log, errs.ErrFmtMore, teamScore.ID, mt.Player)
	}

	if mt.Minute > 90 {
		return nil, nil, 0, errs.ErrGoalMinuteUpperToLimit.Throwf(applog.Log, errs.ErrFmt, mt.Minute)
	}

	return teamScore, player, mt.Minute, nil
}
