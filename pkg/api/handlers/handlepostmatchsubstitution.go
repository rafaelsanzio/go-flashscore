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

func HandlePostMatchSubstitution(w http.ResponseWriter, r *http.Request) {
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

	if match.Status != model.MatchStatusInProgress && match.Status != model.MatchStatusHalftime {
		errs.HttpUnprocessableEntity(w, fmt.Sprintf("err: [%v]", "A sub cannot happen when the game is not in progress or in halftime"))
		return
	}

	matchSubstitutionPayload, err := decodeMatchSubstitutionRequest(r)
	if err != nil {
		errs.HttpUnprocessableEntity(w, fmt.Sprintf("err: [%v]", err.Error()))
		return
	}

	teamInMatch := match.FindTeamInMatch(matchSubstitutionPayload.Team)
	if !teamInMatch {
		errs.HttpUnprocessableEntity(w, fmt.Sprintf("this team is not in this match: [%v]", matchSubstitutionPayload.Team))
		return
	}

	teamSub, playerOut, playerIn, err := convertAndValidatePayloadToMatchSubstitution(ctx, matchSubstitutionPayload)
	if err != nil {
		errs.HttpUnprocessableEntity(w, fmt.Sprintf("err: [%v]", err.Error()))
		return
	}

	minute := matchSubstitutionPayload.Minute
	minuteAsString := strconv.Itoa(minute)

	value := struct {
		TeamScore          team.Team
		PlayerOut          player.Player
		PlayerIn           player.Player
		SubstitutionMinute int
		Created            time.Time
	}{
		TeamScore:          *teamSub,
		PlayerOut:          *playerOut,
		PlayerIn:           *playerIn,
		SubstitutionMinute: minute,
		Created:            time.Now(),
	}

	event := event.Event{
		TournamentID: tournament.ID,
		MatchID:      match.ID,
		Type:         model.EventSubstitution,
		Value:        value,
	}

	err = repo.GetEventRepo().Insert(ctx, event)
	if err != nil {
		errs.HttpInternalServerError(w)
		return
	}

	data := map[string]string{
		"matchEventType":     string(model.EventSubstitution),
		"tournamentID":       tournament.ID,
		"matchID":            match.ID,
		"teamID":             teamSub.ID,
		"playerOutID":        playerOut.ID,
		"playerInID":         playerIn.ID,
		"substitutionMinute": minuteAsString,
	}

	go kafka.Notify(ctx, data, model.ActionGameEvents, "Game Events - Substitution Players", model.KafkaTopicMatchEvents)

	w.WriteHeader(http.StatusCreated)
}

func decodeMatchSubstitutionRequest(r *http.Request) (MatchSubstitutionPayload, errs.AppError) {
	payload := MatchSubstitutionPayload{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		return payload, errs.ErrUnmarshalingJson.Throwf(applog.Log, errs.ErrFmt, err)
	}

	return payload, nil
}

func convertAndValidatePayloadToMatchSubstitution(ctx context.Context, mt MatchSubstitutionPayload) (*team.Team, *player.Player, *player.Player, errs.AppError) {
	team, err := repo.GetTeamRepo().Get(ctx, mt.Team)
	if err != nil || team == nil {
		return nil, nil, nil, errs.ErrTeamIsNotFound.Throwf(applog.Log, errs.ErrFmt, mt.Team)
	}

	if mt.PlayerOut == mt.PlayerIn {
		return nil, nil, nil, errs.ErrSubsSamePlayer.Throwf(applog.Log, errs.ErrFmtMore, mt.PlayerOut, mt.PlayerIn)
	}

	playerOut, err := repo.GetPlayerRepo().GetTeamPlayer(ctx, mt.PlayerOut, team.ID)
	if err != nil || playerOut == nil {
		return nil, nil, nil, errs.ErrPlayerIsNotFoundInThisTeam.Throwf(applog.Log, errs.ErrFmtMore, team.ID, mt.PlayerOut)
	}

	playerIn, err := repo.GetPlayerRepo().GetTeamPlayer(ctx, mt.PlayerIn, team.ID)
	if err != nil || playerIn == nil {
		return nil, nil, nil, errs.ErrPlayerIsNotFoundInThisTeam.Throwf(applog.Log, errs.ErrFmtMore, team.ID, mt.PlayerIn)
	}

	return team, playerOut, playerIn, nil
}
