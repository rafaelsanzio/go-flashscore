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
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
	"github.com/rafaelsanzio/go-flashscore/pkg/prototype"
	"github.com/rafaelsanzio/go-flashscore/pkg/team"
	"github.com/rafaelsanzio/go-flashscore/pkg/tournament"
)

func mockPostTournamentFunc(ctx context.Context, t tournament.Tournament) errs.AppError {
	return nil
}

func mockPostTournamentThrowFunc(ctx context.Context, t tournament.Tournament) errs.AppError {
	return errs.ErrRepoMockAction
}

func TestHandlePostTournament(t *testing.T) {
	body, err := json.Marshal(TournamentEntityPayload{
		Name:  "Any Tournament Name",
		Teams: []string{"any_team_id", "any_team_id_2"},
	})
	assert.Equal(t, nil, err)

	goodReq := httptest.NewRequest(http.MethodPost, "/tournaments", nil)
	goodReq = mux.SetURLVars(goodReq, map[string]string{})
	goodReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	noBodyReq := httptest.NewRequest(http.MethodPost, "/tournaments", nil)
	noBodyReq = mux.SetURLVars(noBodyReq, map[string]string{})

	throwReq := httptest.NewRequest(http.MethodPost, "/tournaments", nil)
	throwReq = mux.SetURLVars(throwReq, map[string]string{})

	throwReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	goodReq2 := httptest.NewRequest(http.MethodPost, "/tournaments", nil)
	goodReq2 = mux.SetURLVars(goodReq2, map[string]string{})
	goodReq2.Body = ioutil.NopCloser(bytes.NewReader(body))

	testCases := []struct {
		Name                     string
		Request                  *http.Request
		HandlePostTournamentFunc func(ctx context.Context, t tournament.Tournament) errs.AppError
		HandleGetTeamFunc        func(ctx context.Context, id string) (*team.Team, errs.AppError)
		ConvertingPayloadFunc    func(ctx context.Context, t TournamentEntityPayload) (*tournament.Tournament, errs.AppError)
		ExpectedStatusCode       int
	}{
		{
			Name:                     "Should return 201 if successful",
			Request:                  goodReq,
			HandlePostTournamentFunc: mockPostTournamentFunc,
			HandleGetTeamFunc:        mockGetTeamFunc,
			ConvertingPayloadFunc:    convertPayloadToTournamentFunc,
			ExpectedStatusCode:       201,
		}, {
			Name:                     "Should return 422 bad request",
			Request:                  noBodyReq,
			HandlePostTournamentFunc: mockPostTournamentFunc,
			HandleGetTeamFunc:        mockGetTeamFunc,
			ConvertingPayloadFunc:    convertPayloadToTournamentFunc,
			ExpectedStatusCode:       422,
		}, {
			Name:                     "Should return 422 throwing error on converting payload func",
			Request:                  goodReq2,
			HandlePostTournamentFunc: mockPostTournamentFunc,
			HandleGetTeamFunc:        mockGetTeamFunc,
			ConvertingPayloadFunc:    fakeConvertPayloadToTournament,
			ExpectedStatusCode:       422,
		}, {
			Name:                     "Should return 500 throwing error on function",
			Request:                  throwReq,
			HandlePostTournamentFunc: mockPostTournamentThrowFunc,
			HandleGetTeamFunc:        mockGetTeamFunc,
			ConvertingPayloadFunc:    convertPayloadToTournamentFunc,
			ExpectedStatusCode:       500,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		repo.SetTeamRepo(repo.MockTeamRepo{
			GetFunc: tc.HandleGetTeamFunc,
		})
		defer repo.SetTeamRepo(nil)

		repo.SetTournamentRepo(repo.MockTournamentRepo{
			InsertFunc: tc.HandlePostTournamentFunc,
		})
		defer repo.SetTournamentRepo(nil)

		convertPayloadToTournamentFunc = tc.ConvertingPayloadFunc
		defer restoreConvertPayloadToTournament(convertPayloadToTournamentFunc)

		w := httptest.NewRecorder()

		HandlePostTournament(w, tc.Request)
		assert.Equal(t, tc.ExpectedStatusCode, w.Code)
	}
}

func TestConvertPayloadToTournament(t *testing.T) {
	inPayload := TournamentEntityPayload{
		Name:  "Any Tournament Name",
		Teams: []string{"any_team_id", "any_team_id_2"},
	}

	expectedTeam := tournament.Tournament{
		ID:    "",
		Name:  "Any Tournament Name",
		Teams: []team.Team{prototype.PrototypeTeam(), prototype.PrototypeTeam()},
	}

	testCases := []struct {
		Name              string
		Payload           TournamentEntityPayload
		HandleGetTeamFunc func(ctx context.Context, id string) (*team.Team, errs.AppError)
		ExpectedTeam      tournament.Tournament
		ExpectError       bool
		ExpectedError     string
	}{
		{
			Name:              "Test Case: 1 - correct body, no error",
			Payload:           inPayload,
			HandleGetTeamFunc: mockGetTeamFunc,
			ExpectedTeam:      expectedTeam,
			ExpectError:       false,
		}, {
			Name:              "Test Case: 2 - throwing error on get function",
			Payload:           inPayload,
			HandleGetTeamFunc: mockGetTeamThrowFunc,
			ExpectedTeam:      expectedTeam,
			ExpectError:       true,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		repo.SetTeamRepo(repo.MockTeamRepo{
			GetFunc: tc.HandleGetTeamFunc,
		})
		defer repo.SetTeamRepo(nil)

		tournament, err := convertPayloadToTournamentFunc(context.Background(), tc.Payload)
		if tc.ExpectError {
			assert.NotNil(t, err)
		} else {
			assert.Equal(t, tc.ExpectedTeam, *tournament)
		}
	}
}

func TestDecodeTournamentRequest(t *testing.T) {
	body, err := json.Marshal(TournamentEntityPayload{
		Name:  "Any Tournament Name",
		Teams: []string{"any_team_id", "any_team_id_2"},
	})
	assert.Equal(t, nil, err)

	goodReq := httptest.NewRequest(http.MethodPost, "/tournaments", nil)
	goodReq = mux.SetURLVars(goodReq, map[string]string{})

	goodReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	noBodyReq := httptest.NewRequest(http.MethodPost, "/tournaments", nil)
	noBodyReq = mux.SetURLVars(noBodyReq, map[string]string{})

	testCases := []struct {
		Name          string
		Request       *http.Request
		Payload       *TournamentEntityPayload
		ExpectedError bool
	}{
		{
			Name:    "Test Case: 1 - correct body, no error",
			Request: goodReq, Payload: &TournamentEntityPayload{
				Name:  "Any Tournament Name",
				Teams: []string{"any_team_id", "any_team_id_2"},
			}, ExpectedError: false,
		}, {
			Name:          "Test Case: 2 - no body, error found",
			Request:       noBodyReq,
			Payload:       nil,
			ExpectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		decodedPayload, err := decodeTournamentRequest(tc.Request)
		if tc.ExpectedError {
			assert.NotNil(t, err)
		} else {
			assert.Equal(t, *tc.Payload, decodedPayload)
		}
	}
}
