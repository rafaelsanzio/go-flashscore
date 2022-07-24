package handlers

import (
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
)

func HandlePostMatchExtratime(w http.ResponseWriter, r *http.Request) {
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
		errs.HttpUnprocessableEntity(w, fmt.Sprintf("err: [%v]", "Match cannot have extra time, because it is not in progress"))
		return
	}

	extratime, err := decodeExtraTimeRequest(r)
	if err != nil {
		errs.HttpUnprocessableEntity(w, fmt.Sprintf("err: [%v]", err.Error()))
		return
	}

	value := struct {
		Extratime int
		Created   time.Time
	}{
		Extratime: extratime,
		Created:   time.Now(),
	}

	event := event.Event{
		TournamentID: tournament.ID,
		MatchID:      match.ID,
		Type:         model.EventExtratime,
		Value:        value,
	}

	extratimeAsString := strconv.Itoa(extratime)

	err = repo.GetEventRepo().Insert(ctx, event)
	if err != nil {
		errs.HttpInternalServerError(w)
		return
	}

	data := map[string]string{
		"matchEventType": string(model.EventExtratime),
		"tournamentID":   tournament.ID,
		"matchID":        match.ID,
		"extratime":      extratimeAsString,
	}

	go kafka.Notify(ctx, data, model.ActionGameEvents, "Game Events - Extratime", model.KafkaTopicMatchEvents)

	w.WriteHeader(http.StatusCreated)
}

func decodeExtraTimeRequest(r *http.Request) (int, errs.AppError) {
	payload := ExtraTimeEntityPayload{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		return 0, errs.ErrUnmarshalingJson.Throwf(applog.Log, errs.ErrFmt, err)
	}

	return payload.Extratime, nil
}
