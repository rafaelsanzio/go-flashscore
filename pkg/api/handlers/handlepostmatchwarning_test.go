package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/event"
	"github.com/rafaelsanzio/go-flashscore/pkg/match"
	"github.com/rafaelsanzio/go-flashscore/pkg/model"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
	"github.com/rafaelsanzio/go-flashscore/pkg/player"
	"github.com/rafaelsanzio/go-flashscore/pkg/prototype"
	"github.com/rafaelsanzio/go-flashscore/pkg/team"
	"github.com/rafaelsanzio/go-flashscore/pkg/tournament"
)

func TestHandlePostMatchWarning(t *testing.T) {
	body, err := json.Marshal(MatchWarningPayload{
		Team:    "1",
		Player:  "1",
		Warning: model.WarningYellowCard,
		Minute:  75,
	})
	assert.Equal(t, nil, err)

	goodReq := httptest.NewRequest(http.MethodPost, "/tournament/{id}/matches/{match_id}/events/warning", nil)
	goodReq = mux.SetURLVars(goodReq, map[string]string{"id": "any", "match_id": "any"})
	goodReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	noBodyReq := httptest.NewRequest(http.MethodPost, "/tournament/{id}/matches/{match_id}/events/warning", nil)
	noBodyReq = mux.SetURLVars(noBodyReq, map[string]string{"id": "any", "match_id": "any"})

	missParamIDReq := httptest.NewRequest(http.MethodPost, "/tournament/{id}/matches/{match_id}/events/warning", nil)
	missParamIDReq = mux.SetURLVars(missParamIDReq, map[string]string{"id": "", "match_id": "any"})
	missParamIDReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	throwReq := httptest.NewRequest(http.MethodPost, "/tournament/{id}/matches/{match_id}/events/warning", nil)
	throwReq = mux.SetURLVars(throwReq, map[string]string{"id": "any", "match_id": "any"})
	throwReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	goodReq2 := httptest.NewRequest(http.MethodPost, "/tournament/{id}/matches/{match_id}/events/warning", nil)
	goodReq2 = mux.SetURLVars(goodReq2, map[string]string{"id": "any", "match_id": "any"})
	goodReq2.Body = ioutil.NopCloser(bytes.NewReader(body))

	missParamMatchIDReq := httptest.NewRequest(http.MethodPost, "/tournament/{id}/matches/{match_id}/events/warning", nil)
	missParamMatchIDReq = mux.SetURLVars(missParamMatchIDReq, map[string]string{"id": "any", "match_id": ""})
	missParamMatchIDReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	bodyMismatchTeams, err := json.Marshal(MatchWarningPayload{
		Team:    "2",
		Player:  "1",
		Warning: model.WarningYellowCard,
		Minute:  70,
	})
	assert.Equal(t, nil, err)

	goodReq3 := httptest.NewRequest(http.MethodPost, "/tournament/{id}/matches/{match_id}/events/substitution", nil)
	goodReq3 = mux.SetURLVars(goodReq3, map[string]string{"id": "any", "match_id": "any"})
	goodReq3.Body = ioutil.NopCloser(bytes.NewReader(bodyMismatchTeams))

	testCases := []struct {
		Name                             string
		Request                          *http.Request
		HandleGetTournamentFunc          func(ctx context.Context, id string) (*tournament.Tournament, errs.AppError)
		HandleFindMatchForTournamentFunc func(ctx context.Context, id, tournamentID string) (*match.Match, errs.AppError)
		HandlePostEventFunc              func(ctx context.Context, e event.Event) errs.AppError
		HandleGetTeamFunc                func(ctx context.Context, id string) (*team.Team, errs.AppError)
		HandleGetTeamPlayerFunc          func(ctx context.Context, id, teamID string) (*player.Player, errs.AppError)
		ExpectedStatusCode               int
	}{
		{
			Name:                             "Should return 201 if successful",
			Request:                          goodReq,
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchStatusInProgressForTournamentFunc,
			HandlePostEventFunc:              mockPostEventFunc,
			HandleGetTeamFunc:                mockGetTeamFunc,
			HandleGetTeamPlayerFunc:          mockGetTeamPlayerSubsFunc,
			ExpectedStatusCode:               201,
		}, {
			Name:                             "Should return 422 if no body request",
			Request:                          noBodyReq,
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchStatusInProgressForTournamentFunc,
			HandlePostEventFunc:              mockPostEventFunc,
			HandleGetTeamFunc:                mockGetTeamFunc,
			HandleGetTeamPlayerFunc:          mockGetTeamPlayerSubsFunc,
			ExpectedStatusCode:               422,
		}, {
			Name:                             "Should return 404 missing id param",
			Request:                          missParamIDReq,
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchStatusInProgressForTournamentFunc,
			HandlePostEventFunc:              mockPostEventFunc,
			HandleGetTeamFunc:                mockGetTeamFunc,
			HandleGetTeamPlayerFunc:          mockGetTeamPlayerSubsFunc,
			ExpectedStatusCode:               404,
		}, {
			Name:                             "Should return 500 throwing error get tournament function",
			Request:                          throwReq,
			HandleGetTournamentFunc:          mockGetTournamentThrowFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchStatusInProgressForTournamentFunc,
			HandlePostEventFunc:              mockPostEventFunc,
			HandleGetTeamFunc:                mockGetTeamFunc,
			HandleGetTeamPlayerFunc:          mockGetTeamPlayerSubsFunc,
			ExpectedStatusCode:               500,
		}, {
			Name:                             "Should return 404 if tournament is not found",
			Request:                          throwReq,
			HandleGetTournamentFunc:          mockGetTournamentNilFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchStatusInProgressForTournamentFunc,
			HandlePostEventFunc:              mockPostEventFunc,
			HandleGetTeamFunc:                mockGetTeamFunc,
			HandleGetTeamPlayerFunc:          mockGetTeamPlayerSubsFunc,
			ExpectedStatusCode:               404,
		}, {
			Name:                             "Should return 404 missing match id param",
			Request:                          missParamMatchIDReq,
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchStatusInProgressForTournamentFunc,
			HandlePostEventFunc:              mockPostEventFunc,
			HandleGetTeamFunc:                mockGetTeamFunc,
			HandleGetTeamPlayerFunc:          mockGetTeamPlayerSubsFunc,
			ExpectedStatusCode:               404,
		}, {
			Name:                             "Should return 500 throwing error find match to tournament function",
			Request:                          throwReq,
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentThrowFunc,
			HandlePostEventFunc:              mockPostEventFunc,
			HandleGetTeamFunc:                mockGetTeamFunc,
			HandleGetTeamPlayerFunc:          mockGetTeamPlayerSubsFunc,
			ExpectedStatusCode:               500,
		}, {
			Name:                             "Should return 404 if match is not found",
			Request:                          throwReq,
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentNilFunc,
			HandlePostEventFunc:              mockPostEventFunc,
			HandleGetTeamFunc:                mockGetTeamFunc,
			HandleGetTeamPlayerFunc:          mockGetTeamPlayerSubsFunc,
			ExpectedStatusCode:               404,
		}, {
			Name:                             "Should return 422 if match is not started",
			Request:                          throwReq,
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchStatusNotStartedForTournamentFunc,
			HandlePostEventFunc:              mockPostEventFunc,
			HandleGetTeamFunc:                mockGetTeamFunc,
			HandleGetTeamPlayerFunc:          mockGetTeamPlayerSubsFunc,
			ExpectedStatusCode:               422,
		}, {
			Name:                             "Should return 500 throwing error post event function",
			Request:                          throwReq,
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchStatusInProgressForTournamentFunc,
			HandlePostEventFunc:              mockPostEventThrowFunc,
			HandleGetTeamFunc:                mockGetTeamFunc,
			HandleGetTeamPlayerFunc:          mockGetTeamPlayerSubsFunc,
			ExpectedStatusCode:               500,
		}, {
			Name:                             "Should return 422 if if the teams is not in the match",
			Request:                          goodReq3,
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchStatusInProgressForTournamentFunc,
			HandlePostEventFunc:              mockPostEventFunc,
			HandleGetTeamFunc:                mockGetTeamFunc,
			HandleGetTeamPlayerFunc:          mockGetTeamPlayerSubsFunc,
			ExpectedStatusCode:               422,
		}, {
			Name:                             "Should return 422 throwing error get team function",
			Request:                          goodReq2,
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchStatusInProgressForTournamentFunc,
			HandlePostEventFunc:              mockPostEventFunc,
			HandleGetTeamFunc:                mockGetTeamThrowFunc,
			HandleGetTeamPlayerFunc:          mockGetTeamPlayerSubsFunc,
			ExpectedStatusCode:               422,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		repo.SetTournamentRepo(repo.MockTournamentRepo{
			GetFunc: tc.HandleGetTournamentFunc,
		})
		defer repo.SetTournamentRepo(nil)

		repo.SetMatchRepo(repo.MockMatchRepo{
			FindMatchForTournamentFunc: tc.HandleFindMatchForTournamentFunc,
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

		repo.SetEventRepo(repo.MockEventRepo{
			InsertFunc: tc.HandlePostEventFunc,
		})
		defer repo.SetEventRepo(nil)

		w := httptest.NewRecorder()

		HandlePostMatchWarning(w, tc.Request)
		assert.Equal(t, tc.ExpectedStatusCode, w.Code)
	}
}

type payloadWarningReturn struct {
	team    team.Team
	player  player.Player
	warning model.Warnings
	minute  int
}

func TestConvertAndValidatePayloadToMatchWarning(t *testing.T) {
	inPayload := MatchWarningPayload{
		Team:    "1",
		Player:  "1",
		Warning: "RedCard",
		Minute:  70,
	}

	expectPayload := payloadWarningReturn{
		team:    prototype.PrototypeTeam(),
		player:  prototype.PrototypePlayer(),
		warning: model.WarningRedCard,
		minute:  70,
	}

	testCases := []struct {
		Name                    string
		Payload                 MatchWarningPayload
		HandleGetTeamFunc       func(ctx context.Context, id string) (*team.Team, errs.AppError)
		HandleGetTeamPlayerFunc func(ctx context.Context, id, teamID string) (*player.Player, errs.AppError)
		ExpectedReturn          payloadWarningReturn
		ExpectError             bool
	}{
		{
			Name:                    "Test Case: 1 - correct body, no error",
			Payload:                 inPayload,
			HandleGetTeamFunc:       mockGetTeamFuncForMatch,
			HandleGetTeamPlayerFunc: mockGetTeamPlayerSubsFunc,
			ExpectedReturn:          expectPayload,
			ExpectError:             false,
		}, {
			Name:                    "Test Case: 2 - throwing error on get team function",
			Payload:                 inPayload,
			HandleGetTeamFunc:       mockGetTeamThrowFunc,
			HandleGetTeamPlayerFunc: mockGetTeamPlayerSubsFunc,
			ExpectedReturn:          expectPayload,
			ExpectError:             true,
		}, {
			Name:                    "Test Case: 3 - throwing error get player function",
			Payload:                 inPayload,
			HandleGetTeamFunc:       mockGetTeamFuncForMatch,
			HandleGetTeamPlayerFunc: mockGetTeamPlayerThrowFunc,
			ExpectedReturn:          expectPayload,
			ExpectError:             true,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		repo.SetTeamRepo(repo.MockTeamRepo{
			GetFunc: tc.HandleGetTeamFunc,
		})
		defer repo.SetTeamRepo(nil)

		repo.SetPlayerRepo(repo.MockPlayerRepo{
			GetTeamPlayerFunc: tc.HandleGetTeamPlayerFunc,
		})
		defer repo.SetPlayerRepo(nil)

		team, player, err := convertAndValidatePayloadToMatchWarning(context.Background(), tc.Payload)
		if tc.ExpectError {
			assert.NotNil(t, err)
		} else {
			assert.Equal(t, tc.ExpectedReturn.team, *team)
			assert.Equal(t, tc.ExpectedReturn.player, *player)
		}
	}
}
