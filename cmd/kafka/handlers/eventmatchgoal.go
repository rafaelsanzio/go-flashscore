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

func HandleEventMatchGoal(ctx context.Context, data map[string]string) errs.AppError {
	tournamentID := data["tournamentID"]
	matchID := data["matchID"]
	teamScoreID := data["teamScore"]
	playerScoreID := data["player"]
	goalMinute := data["goalMinute"]

	tournament, err := repo.GetTournamentRepo().Get(ctx, tournamentID)
	if err != nil || tournament == nil {
		return errs.ErrTournamentIsNotFound.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}

	match, err := repo.GetMatchRepo().FindMatchForTournament(ctx, matchID, tournamentID)
	if err != nil || match == nil {
		return errs.ErrMatchIsNotFound.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}

	teamScore, err := repo.GetTeamRepo().Get(ctx, teamScoreID)
	if err != nil || teamScore == nil {
		return errs.ErrTeamIsNotFound.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}

	playerScore, err := repo.GetPlayerRepo().GetTeamPlayer(ctx, playerScoreID, teamScore.ID)
	if err != nil || playerScore == nil {
		return errs.ErrPlayerIsNotFoundInThisTeam.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}

	goalMinuteAsInt, err_ := strconvAtoi(goalMinute)
	if err_ != nil {
		return errs.ErrParsingAtoi.Throwf(applog.Log, errs.ErrFmt, err_.Error())
	}

	events := struct {
		MatchEvent model.EventsMatchType
		TeamScore  team.Team
		Player     player.Player
		GoalMinute int
		Created    time.Time
	}{
		MatchEvent: model.EventGoal,
		TeamScore:  *teamScore,
		Player:     *playerScore,
		GoalMinute: goalMinuteAsInt,
		Created:    time.Now(),
	}
	match.Events = append(match.Events, events)

	matchUpdated, err := repo.GetMatchRepo().Update(ctx, *match)
	if err != nil || matchUpdated == nil {
		return errs.ErrMatchIsNotFound.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}

	return nil
}
