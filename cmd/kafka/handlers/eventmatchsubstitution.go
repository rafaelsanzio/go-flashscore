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

func HandleEventMatchSubstitution(ctx context.Context, data map[string]string) errs.AppError {
	tournamentID := data["tournamentID"]
	matchID := data["matchID"]
	teamID := data["teamID"]
	playerOutID := data["playerOutID"]
	playerInID := data["playerInID"]
	substitutionMinute := data["substitutionMinute"]

	tournament, err := repo.GetTournamentRepo().Get(ctx, tournamentID)
	if err != nil || tournament == nil {
		return errs.ErrTournamentIsNotFound.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}

	match, err := repo.GetMatchRepo().FindMatchForTournament(ctx, matchID, tournamentID)
	if err != nil || match == nil {
		return errs.ErrMatchIsNotFound.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}

	teamSub, err := repo.GetTeamRepo().Get(ctx, teamID)
	if err != nil || teamSub == nil {
		return errs.ErrTeamIsNotFound.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}

	playerOut, err := repo.GetPlayerRepo().GetTeamPlayer(ctx, playerOutID, teamID)
	if err != nil || playerOut == nil {
		return errs.ErrPlayerIsNotFoundInThisTeam.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}

	playerIn, err := repo.GetPlayerRepo().GetTeamPlayer(ctx, playerInID, teamID)
	if err != nil || playerIn == nil {
		return errs.ErrPlayerIsNotFoundInThisTeam.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}

	substitutionMinuteAsInt, err_ := strconvAtoi(substitutionMinute)
	if err_ != nil {
		return errs.ErrParsingAtoi.Throwf(applog.Log, errs.ErrFmt, err_.Error())
	}

	events := struct {
		MatchEvent         model.EventsMatchType
		Team               team.Team
		PlayerOut          player.Player
		PlayerIn           player.Player
		SubstitutionMinute int
		Created            time.Time
	}{
		MatchEvent:         model.EventSubstitution,
		Team:               *teamSub,
		PlayerOut:          *playerOut,
		PlayerIn:           *playerIn,
		SubstitutionMinute: substitutionMinuteAsInt,
		Created:            time.Now(),
	}
	match.Events = append(match.Events, events)

	matchUpdated, err := repo.GetMatchRepo().Update(ctx, *match)
	if err != nil || matchUpdated == nil {
		return errs.ErrMatchIsNotFound.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}

	return nil
}
