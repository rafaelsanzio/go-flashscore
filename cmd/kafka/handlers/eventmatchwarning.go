package handlers

import (
	"context"
	"time"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/model"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
	"github.com/rafaelsanzio/go-flashscore/pkg/player"
	"github.com/rafaelsanzio/go-flashscore/pkg/team"
)

func HandleEventMatchWarning(ctx context.Context, data map[string]string) errs.AppError {
	tournamentID := data["tournamentID"]
	matchID := data["matchID"]
	teamID := data["teamID"]
	playerID := data["playerID"]
	warning := data["warning"]
	warningMinute := data["warningMinute"]

	tournament, err := repo.GetTournamentRepo().Get(ctx, tournamentID)
	if err != nil || tournament == nil {
		return errs.ErrTournamentIsNotFound.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}

	match, err := repo.GetMatchRepo().FindMatchForTournament(ctx, matchID, tournamentID)
	if err != nil || match == nil {
		return errs.ErrMatchIsNotFound.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}

	teamWarn, err := repo.GetTeamRepo().Get(ctx, teamID)
	if err != nil || teamWarn == nil {
		return errs.ErrTeamIsNotFound.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}

	playerWarn, err := repo.GetPlayerRepo().GetTeamPlayer(ctx, playerID, teamID)
	if err != nil || playerWarn == nil {
		return errs.ErrPlayerIsNotFoundInThisTeam.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}

	warningMinuteAsInt, err_ := strconvAtoi(warningMinute)
	if err_ != nil {
		return errs.ErrParsingAtoi.Throwf(applog.Log, errs.ErrFmt, err_.Error())
	}

	events := struct {
		MatchEvent    model.EventsMatchType
		Team          team.Team
		Player        player.Player
		Warning       model.Warnings
		WarningMinute int
		Created       time.Time
	}{
		MatchEvent:    model.EventWarning,
		Team:          *teamWarn,
		Player:        *playerWarn,
		Warning:       model.Warnings(warning),
		WarningMinute: warningMinuteAsInt,
		Created:       time.Now(),
	}
	match.Events = append(match.Events, events)

	matchUpdated, err := repo.GetMatchRepo().Update(ctx, *match)
	if err != nil || matchUpdated == nil {
		return errs.ErrMatchIsNotFound.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}

	return nil
}
