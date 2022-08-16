package notification

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/match"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
	"github.com/rafaelsanzio/go-flashscore/pkg/player"
	"github.com/rafaelsanzio/go-flashscore/pkg/prototype"
	"github.com/rafaelsanzio/go-flashscore/pkg/team"
	"github.com/rafaelsanzio/go-flashscore/pkg/tournament"
)

func mockGetPlayerFunc(ctx context.Context, id string) (*player.Player, errs.AppError) {
	playerMock := prototype.PrototypePlayer()
	return &playerMock, nil
}

func mockGetPlayerThrowFunc(ctx context.Context, id string) (*player.Player, errs.AppError) {
	return nil, errs.ErrRepoMockAction
}

func mockUpdatePlayerFunc(ctx context.Context, p player.Player) (*player.Player, errs.AppError) {
	return &p, nil
}

func mockGetTeamFunc(ctx context.Context, id string) (*team.Team, errs.AppError) {
	teamMock := prototype.PrototypeTeam()
	return &teamMock, nil
}

func TestHandlerUpdateTeamPlayer(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		Name                   string
		Body                   string
		HandleGetPlayerFunc    func(ctx context.Context, id string) (*player.Player, errs.AppError)
		HandleGetTeamFunc      func(ctx context.Context, id string) (*team.Team, errs.AppError)
		HandleUpdatePlayerFunc func(ctx context.Context, t player.Player) (*player.Player, errs.AppError)
		ExpectedError          bool
	}{
		{
			Name:                   "Handle update team player correct",
			Body:                   `{"Action":"UpdateTeamPlayer", "Data":{"playerID":"any-player-id", "teamDestinyID":"any-team-id"}}`,
			HandleGetPlayerFunc:    mockGetPlayerFunc,
			HandleGetTeamFunc:      mockGetTeamFunc,
			HandleUpdatePlayerFunc: mockUpdatePlayerFunc,
			ExpectedError:          false,
		}, {
			Name:                   "Handle update team player throw error on get player function",
			Body:                   `{"Action":"UpdateTeamPlayer", "Data":{"playerID":"any-player-id", "teamDestinyID":"any-team-id"}}`,
			HandleGetPlayerFunc:    mockGetPlayerThrowFunc,
			HandleGetTeamFunc:      mockGetTeamFunc,
			HandleUpdatePlayerFunc: mockUpdatePlayerFunc,
			ExpectedError:          true,
		}, {
			Name:                   "Handle update team player throw error on unmarshal function",
			Body:                   `{"Action":"UpdateTeamPlayer", "Data":{"playerID":"any-player-id", "teamDestinyID":"any-team-id"}`,
			HandleGetPlayerFunc:    mockGetPlayerFunc,
			HandleGetTeamFunc:      mockGetTeamFunc,
			HandleUpdatePlayerFunc: mockUpdatePlayerFunc,
			ExpectedError:          true,
		}, {
			Name:                   "Handle update team player throw error on non exist action",
			Body:                   `{"Action":"Unknown", "Data":{"playerID":"any-player-id", "teamDestinyID":"any-team-id"}}`,
			HandleGetPlayerFunc:    mockGetPlayerFunc,
			HandleGetTeamFunc:      mockGetTeamFunc,
			HandleUpdatePlayerFunc: mockUpdatePlayerFunc,
			ExpectedError:          true,
		},
	}

	for _, tc := range testCases {
		t.Logf(tc.Name)

		repo.SetPlayerRepo(repo.MockPlayerRepo{
			GetFunc:    tc.HandleGetPlayerFunc,
			UpdateFunc: tc.HandleUpdatePlayerFunc,
		})
		defer repo.SetPlayerRepo(nil)

		repo.SetTeamRepo(repo.MockTeamRepo{
			GetFunc: tc.HandleGetTeamFunc,
		})
		defer repo.SetTeamRepo(nil)

		err := Handler(ctx, tc.Body, "any-key")
		if tc.ExpectedError {
			assert.NotNil(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func mockGetTournamentFunc(ctx context.Context, id string) (*tournament.Tournament, errs.AppError) {
	tournamentMock := prototype.PrototypeTournament()
	return &tournamentMock, nil
}

func mockGetTournamentThrowFunc(ctx context.Context, id string) (*tournament.Tournament, errs.AppError) {
	return nil, errs.ErrRepoMockAction
}

func mockFindMatchForTournamentFunc(ctx context.Context, id string, tournamentID string) (*match.Match, errs.AppError) {
	matchMock := prototype.PrototypeMatch()
	return &matchMock, nil
}

func mockUpdateMatchFunc(ctx context.Context, m match.Match) (*match.Match, errs.AppError) {
	return &m, nil
}

func TestHandlerMatchEventStart(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		Name                       string
		Body                       string
		HandleGetTournamentFunc    func(ctx context.Context, id string) (*tournament.Tournament, errs.AppError)
		FindMatchForTournamentFunc func(ctx context.Context, id string, tournamentID string) (*match.Match, errs.AppError)
		UpdateMatchFunc            func(ctx context.Context, m match.Match) (*match.Match, errs.AppError)
		ExpectedError              bool
	}{
		{
			Name:                       "Handle action game event match start",
			Body:                       `{"Action":"ActionGameEvents","Data":{"matchEventType":"Start", "tournamentID":"any-tournament-id", "matchID":"any-match-id", "timeStarted":"16:00"}}`,
			HandleGetTournamentFunc:    mockGetTournamentFunc,
			FindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			UpdateMatchFunc:            mockUpdateMatchFunc,
			ExpectedError:              false,
		}, {
			Name:                       "Handle action game event match start error",
			Body:                       `{"Action":"ActionGameEvents","Data":{"matchEventType":"Start", "tournamentID":"any-tournament-id", "matchID":"any-match-id", "timeStarted":"16:00"}}`,
			HandleGetTournamentFunc:    mockGetTournamentThrowFunc,
			FindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			UpdateMatchFunc:            mockUpdateMatchFunc,
			ExpectedError:              true,
		},
	}

	for _, tc := range testCases {
		t.Logf(tc.Name)

		repo.SetTournamentRepo(repo.MockTournamentRepo{
			GetFunc: tc.HandleGetTournamentFunc,
		})
		defer repo.SetTournamentRepo(nil)

		repo.SetMatchRepo(repo.MockMatchRepo{
			FindMatchForTournamentFunc: tc.FindMatchForTournamentFunc,
			UpdateFunc:                 tc.UpdateMatchFunc,
		})
		defer repo.SetMatchRepo(nil)

		err := Handler(ctx, tc.Body, "any-key")
		if tc.ExpectedError {
			assert.NotNil(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func mockGetTeamPlayerFunc(ctx context.Context, id string, teamID string) (*player.Player, errs.AppError) {
	playerMock := prototype.PrototypePlayer()
	return &playerMock, nil
}

func TestHandlerMatchEventGoal(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		Name                       string
		Body                       string
		HandleGetTournamentFunc    func(ctx context.Context, id string) (*tournament.Tournament, errs.AppError)
		FindMatchForTournamentFunc func(ctx context.Context, id string, tournamentID string) (*match.Match, errs.AppError)
		HandleGetTeamFunc          func(ctx context.Context, id string) (*team.Team, errs.AppError)
		HandleGetTeamPlayerFunc    func(ctx context.Context, id string, teamID string) (*player.Player, errs.AppError)
		UpdateMatchFunc            func(ctx context.Context, m match.Match) (*match.Match, errs.AppError)
		ExpectedError              bool
	}{
		{
			Name:                       "Handle action game event match goal",
			Body:                       `{"Action":"ActionGameEvents","Data":{"matchEventType":"Goal", "tournamentID":"any-tournament-id", "matchID":"any-match-id", "player":"any-player-id", "goalMinute":"10"}}`,
			HandleGetTournamentFunc:    mockGetTournamentFunc,
			FindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			HandleGetTeamFunc:          mockGetTeamFunc,
			HandleGetTeamPlayerFunc:    mockGetTeamPlayerFunc,
			UpdateMatchFunc:            mockUpdateMatchFunc,
			ExpectedError:              false,
		}, {
			Name:                       "Handle action game event match goal error",
			Body:                       `{"Action":"ActionGameEvents","Data":{"matchEventType":"Goal", "tournamentID":"any-tournament-id", "matchID":"any-match-id", "player":"any-player-id", "goalMinute":"10"}}`,
			HandleGetTournamentFunc:    mockGetTournamentThrowFunc,
			FindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			HandleGetTeamFunc:          mockGetTeamFunc,
			HandleGetTeamPlayerFunc:    mockGetTeamPlayerFunc,
			UpdateMatchFunc:            mockUpdateMatchFunc,
			ExpectedError:              true,
		},
	}

	for _, tc := range testCases {
		t.Logf(tc.Name)

		repo.SetTournamentRepo(repo.MockTournamentRepo{
			GetFunc: tc.HandleGetTournamentFunc,
		})
		defer repo.SetTournamentRepo(nil)

		repo.SetTeamRepo(repo.MockTeamRepo{
			GetFunc: tc.HandleGetTeamFunc,
		})
		defer repo.SetTeamRepo(nil)

		repo.SetPlayerRepo(repo.MockPlayerRepo{
			GetTeamPlayerFunc: tc.HandleGetTeamPlayerFunc,
		})
		defer repo.SetPlayerRepo(nil)

		repo.SetMatchRepo(repo.MockMatchRepo{
			FindMatchForTournamentFunc: tc.FindMatchForTournamentFunc,
			UpdateFunc:                 tc.UpdateMatchFunc,
		})
		defer repo.SetMatchRepo(nil)

		err := Handler(ctx, tc.Body, "any-key")
		if tc.ExpectedError {
			assert.NotNil(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestHandlerMatchEventHalftime(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		Name                       string
		Body                       string
		HandleGetTournamentFunc    func(ctx context.Context, id string) (*tournament.Tournament, errs.AppError)
		FindMatchForTournamentFunc func(ctx context.Context, id string, tournamentID string) (*match.Match, errs.AppError)
		UpdateMatchFunc            func(ctx context.Context, m match.Match) (*match.Match, errs.AppError)
		ExpectedError              bool
	}{
		{
			Name:                       "Handle action game event match halftime",
			Body:                       `{"Action":"ActionGameEvents","Data":{"matchEventType":"Halftime", "tournamentID":"any-tournament-id","matchID":"any-match-id", "halftime":"any-halftime"}}`,
			HandleGetTournamentFunc:    mockGetTournamentFunc,
			FindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			UpdateMatchFunc:            mockUpdateMatchFunc,
			ExpectedError:              false,
		}, {
			Name:                       "Handle action game event match halftime error",
			Body:                       `{"Action":"ActionGameEvents","Data":{"matchEventType":"Halftime", "tournamentID":"any-tournament-id","matchID":"any-match-id", "halftime":"any-halftime"}}`,
			HandleGetTournamentFunc:    mockGetTournamentThrowFunc,
			FindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			UpdateMatchFunc:            mockUpdateMatchFunc,
			ExpectedError:              true,
		},
	}

	for _, tc := range testCases {
		t.Logf(tc.Name)

		repo.SetTournamentRepo(repo.MockTournamentRepo{
			GetFunc: tc.HandleGetTournamentFunc,
		})
		defer repo.SetTournamentRepo(nil)

		repo.SetMatchRepo(repo.MockMatchRepo{
			FindMatchForTournamentFunc: tc.FindMatchForTournamentFunc,
			UpdateFunc:                 tc.UpdateMatchFunc,
		})
		defer repo.SetMatchRepo(nil)

		err := Handler(ctx, tc.Body, "any-key")
		if tc.ExpectedError {
			assert.NotNil(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestHandlerMatchEventNotFound(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		Name          string
		Body          string
		ExpectedError bool
	}{
		{
			Name:          "Handle action game event that don't exist",
			Body:          `{"Action":"ActionGameEvents","Data":{"matchEventType":"Default", "tournamentID":"any-tournament-id", "matchID":"any-match-id"}}`,
			ExpectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Logf(tc.Name)

		err := Handler(ctx, tc.Body, "any-key")
		if tc.ExpectedError {
			assert.NotNil(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestHandlerMatchEventExtratime(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		Name                       string
		Body                       string
		HandleGetTournamentFunc    func(ctx context.Context, id string) (*tournament.Tournament, errs.AppError)
		FindMatchForTournamentFunc func(ctx context.Context, id string, tournamentID string) (*match.Match, errs.AppError)
		UpdateMatchFunc            func(ctx context.Context, m match.Match) (*match.Match, errs.AppError)
		ExpectedError              bool
	}{
		{
			Name:                       "Handle action game event match extratime",
			Body:                       `{"Action":"ActionGameEvents","Data":{"matchEventType":"Extratime", "tournamentID":"any-tournament-id", "matchID":"any-match-id", "extratime":"5"}}`,
			HandleGetTournamentFunc:    mockGetTournamentFunc,
			FindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			UpdateMatchFunc:            mockUpdateMatchFunc,
			ExpectedError:              false,
		}, {
			Name:                       "Handle action game event match extratime error",
			Body:                       `{"Action":"ActionGameEvents","Data":{"matchEventType":"Extratime", "tournamentID":"any-tournament-id", "matchID":"any-match-id", "extratime":"5"}}`,
			HandleGetTournamentFunc:    mockGetTournamentThrowFunc,
			FindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			UpdateMatchFunc:            mockUpdateMatchFunc,
			ExpectedError:              true,
		},
	}

	for _, tc := range testCases {
		t.Logf(tc.Name)

		repo.SetTournamentRepo(repo.MockTournamentRepo{
			GetFunc: tc.HandleGetTournamentFunc,
		})
		defer repo.SetTournamentRepo(nil)

		repo.SetMatchRepo(repo.MockMatchRepo{
			FindMatchForTournamentFunc: tc.FindMatchForTournamentFunc,
			UpdateFunc:                 tc.UpdateMatchFunc,
		})
		defer repo.SetMatchRepo(nil)

		err := Handler(ctx, tc.Body, "any-key")
		if tc.ExpectedError {
			assert.NotNil(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestHandlerMatchEventSubstitution(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		Name                       string
		Body                       string
		HandleGetTournamentFunc    func(ctx context.Context, id string) (*tournament.Tournament, errs.AppError)
		FindMatchForTournamentFunc func(ctx context.Context, id string, tournamentID string) (*match.Match, errs.AppError)
		HandleGetTeamFunc          func(ctx context.Context, id string) (*team.Team, errs.AppError)
		HandleGetTeamPlayerFunc    func(ctx context.Context, id string, teamID string) (*player.Player, errs.AppError)
		UpdateMatchFunc            func(ctx context.Context, m match.Match) (*match.Match, errs.AppError)
		ExpectedError              bool
	}{
		{
			Name:                       "Handle action game event match substitution",
			Body:                       `{"Action":"ActionGameEvents","Data":{"matchEventType":"Substitution", "tournamentID":"any-tournament-id", "matchID":"any-match-id", "teamID": "any-team-id", "playerOutID": "any-player-out-id", "playerInID": "any-player-in-id", "substitutionMinute":"5"}}`,
			HandleGetTournamentFunc:    mockGetTournamentFunc,
			FindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			HandleGetTeamFunc:          mockGetTeamFunc,
			HandleGetTeamPlayerFunc:    mockGetTeamPlayerFunc,
			UpdateMatchFunc:            mockUpdateMatchFunc,
			ExpectedError:              false,
		}, {
			Name:                       "Handle action game event match substitution error",
			Body:                       `{"Action":"ActionGameEvents","Data":{"matchEventType":"Substitution", "tournamentID":"any-tournament-id", "matchID":"any-match-id", "teamID": "any-team-id", "playerOutID": "any-player-out-id", "playerInID": "any-player-in-id", "substitutionMinute":"5"}}`,
			HandleGetTournamentFunc:    mockGetTournamentThrowFunc,
			FindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			HandleGetTeamFunc:          mockGetTeamFunc,
			HandleGetTeamPlayerFunc:    mockGetTeamPlayerFunc,
			UpdateMatchFunc:            mockUpdateMatchFunc,
			ExpectedError:              true,
		},
	}

	for _, tc := range testCases {
		t.Logf(tc.Name)

		repo.SetTournamentRepo(repo.MockTournamentRepo{
			GetFunc: tc.HandleGetTournamentFunc,
		})
		defer repo.SetTournamentRepo(nil)

		repo.SetMatchRepo(repo.MockMatchRepo{
			FindMatchForTournamentFunc: tc.FindMatchForTournamentFunc,
			UpdateFunc:                 tc.UpdateMatchFunc,
		})
		defer repo.SetMatchRepo(nil)

		repo.SetTeamRepo(repo.MockTeamRepo{
			GetFunc: tc.HandleGetTeamFunc,
		})
		defer repo.SetTeamRepo(nil)

		repo.SetPlayerRepo(repo.MockPlayerRepo{
			GetTeamPlayerFunc: tc.HandleGetTeamPlayerFunc,
		})
		defer repo.SetPlayerRepo(nil)

		err := Handler(ctx, tc.Body, "any-key")
		if tc.ExpectedError {
			assert.NotNil(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestHandlerMatchEventWarning(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		Name                       string
		Body                       string
		HandleGetTournamentFunc    func(ctx context.Context, id string) (*tournament.Tournament, errs.AppError)
		FindMatchForTournamentFunc func(ctx context.Context, id string, tournamentID string) (*match.Match, errs.AppError)
		HandleGetTeamFunc          func(ctx context.Context, id string) (*team.Team, errs.AppError)
		HandleGetTeamPlayerFunc    func(ctx context.Context, id string, teamID string) (*player.Player, errs.AppError)
		UpdateMatchFunc            func(ctx context.Context, m match.Match) (*match.Match, errs.AppError)
		ExpectedError              bool
	}{
		{
			Name:                       "Handle action game event match warning",
			Body:                       `{"Action":"ActionGameEvents","Data":{"matchEventType":"Warning", "tournamentID":"any-tournament-id", "matchID":"any-match-id", "teamID": "any-team-id", "playerID": "any-player-id", "warning": "any-warning", "warningMinute":"5"}}`,
			HandleGetTournamentFunc:    mockGetTournamentFunc,
			FindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			HandleGetTeamFunc:          mockGetTeamFunc,
			HandleGetTeamPlayerFunc:    mockGetTeamPlayerFunc,
			UpdateMatchFunc:            mockUpdateMatchFunc,
			ExpectedError:              false,
		}, {
			Name:                       "Handle action game event match warning error",
			Body:                       `{"Action":"ActionGameEvents","Data":{"matchEventType":"Warning", "tournamentID":"any-tournament-id", "matchID":"any-match-id", "teamID": "any-team-id", "playerID": "any-player-id", "warning": "any-warning", "warningMinute":"5"}}`,
			HandleGetTournamentFunc:    mockGetTournamentThrowFunc,
			FindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			HandleGetTeamFunc:          mockGetTeamFunc,
			HandleGetTeamPlayerFunc:    mockGetTeamPlayerFunc,
			UpdateMatchFunc:            mockUpdateMatchFunc,
			ExpectedError:              true,
		},
	}

	for _, tc := range testCases {
		t.Logf(tc.Name)

		repo.SetTournamentRepo(repo.MockTournamentRepo{
			GetFunc: tc.HandleGetTournamentFunc,
		})
		defer repo.SetTournamentRepo(nil)

		repo.SetMatchRepo(repo.MockMatchRepo{
			FindMatchForTournamentFunc: tc.FindMatchForTournamentFunc,
			UpdateFunc:                 tc.UpdateMatchFunc,
		})
		defer repo.SetMatchRepo(nil)

		repo.SetTeamRepo(repo.MockTeamRepo{
			GetFunc: tc.HandleGetTeamFunc,
		})
		defer repo.SetTeamRepo(nil)

		repo.SetPlayerRepo(repo.MockPlayerRepo{
			GetTeamPlayerFunc: tc.HandleGetTeamPlayerFunc,
		})
		defer repo.SetPlayerRepo(nil)

		err := Handler(ctx, tc.Body, "any-key")
		if tc.ExpectedError {
			assert.NotNil(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
