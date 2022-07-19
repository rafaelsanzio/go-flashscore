package handlers

import (
	"context"
	"time"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/model"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
)

func HandleEventMatchExtratime(ctx context.Context, data map[string]string) errs.AppError {
	tournamentID := data["tournamentID"]
	matchID := data["matchID"]
	extratime := data["extratime"]

	tournament, err := repo.GetTournamentRepo().Get(ctx, tournamentID)
	if err != nil || tournament == nil {
		return errs.ErrTournamentIsNotFound.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}

	match, err := repo.GetMatchRepo().FindMatchForTournament(ctx, matchID, tournamentID)
	if err != nil || match == nil {
		return errs.ErrMatchIsNotFound.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}

	extratimeAsInt, err_ := strconvAtoi(extratime)
	if err_ != nil {
		return errs.ErrParsingAtoi.Throwf(applog.Log, errs.ErrFmt, err_.Error())
	}

	events := struct {
		MatchEvent model.EventsMatchType
		Extratime  int
		Created    time.Time
	}{
		MatchEvent: model.EventExtratime,
		Extratime:  extratimeAsInt,
		Created:    time.Now(),
	}
	match.Events = append(match.Events, events)

	matchUpdated, err := repo.GetMatchRepo().Update(ctx, *match)
	if err != nil || matchUpdated == nil {
		return errs.ErrMatchIsNotFound.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}

	return nil
}
