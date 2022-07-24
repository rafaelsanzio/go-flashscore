package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/event"
	"github.com/rafaelsanzio/go-flashscore/pkg/kafka"
	"github.com/rafaelsanzio/go-flashscore/pkg/model"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
)

func HandlePostMatchFinish(w http.ResponseWriter, r *http.Request) {
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
		errs.HttpUnprocessableEntity(w, fmt.Sprintf("err: [%v]", "Match is not in progress to be finished"))
		return
	}

	timeFinished := time.Now().Format("15:04")

	value := struct {
		TimeFinished string
		Created      time.Time
	}{
		TimeFinished: timeFinished,
		Created:      time.Now(),
	}

	event := event.Event{
		TournamentID: tournament.ID,
		MatchID:      match.ID,
		Type:         model.EventFinish,
		Value:        value,
	}

	err = repo.GetEventRepo().Insert(ctx, event)
	if err != nil {
		errs.HttpInternalServerError(w)
		return
	}

	data := map[string]string{
		"matchEventType": string(model.EventFinish),
		"tournamentID":   tournament.ID,
		"matchID":        match.ID,
		"timeFinished":   timeFinished,
	}

	go kafka.Notify(ctx, data, model.ActionGameEvents, "Game Events - Finish", model.KafkaTopicMatchEvents)

	w.WriteHeader(http.StatusCreated)
}
