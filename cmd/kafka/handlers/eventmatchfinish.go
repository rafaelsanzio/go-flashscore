package handlers

import (
	"context"
	"time"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/model"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
)

func HandleEventMatchFinish(ctx context.Context, data map[string]string) errs.AppError {
	tournamentID := data["tournamentID"]
	matchID := data["matchID"]
	timeFinished := data["timeFinished"]

	tournament, err := repo.GetTournamentRepo().Get(ctx, tournamentID)
	if err != nil || tournament == nil {
		return errs.ErrTournamentIsNotFound.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}

	match, err := repo.GetMatchRepo().FindMatchForTournament(ctx, matchID, tournamentID)
	if err != nil || match == nil {
		return errs.ErrMatchIsNotFound.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}

	events := struct {
		MatchEvent   model.EventsMatchType
		TimeFinished string
		Created      time.Time
	}{
		MatchEvent:   model.EventFinish,
		TimeFinished: timeFinished,
		Created:      time.Now(),
	}
	match.Events = append(match.Events, events)

	match.Status = model.MatchStatusFinished

	matchUpdated, err := repo.GetMatchRepo().Update(ctx, *match)
	if err != nil || matchUpdated == nil {
		return errs.ErrMatchIsNotFound.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}

	return nil
}
