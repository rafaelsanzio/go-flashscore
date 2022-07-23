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
	"github.com/rafaelsanzio/go-flashscore/pkg/team"
	"github.com/rafaelsanzio/go-flashscore/pkg/tournament"
)

func TestHandleAddTeamsTournament(t *testing.T) {
	body, err := json.Marshal(AddTeamsTournamentEntityPayload{
		Teams: []string{"any_team_id", "any_team_id_2"},
	})
	assert.Equal(t, nil, err)

	emptyBody, err := json.Marshal(AddTeamsTournamentEntityPayload{
		Teams: []string{},
	})
	assert.Equal(t, nil, err)

	goodReq := httptest.NewRequest(http.MethodPost, "/tournaments/:id/add-teams", nil)
	goodReq = mux.SetURLVars(goodReq, map[string]string{"id": "any"})
	goodReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	noBodyReq := httptest.NewRequest(http.MethodPost, "/tournaments/:id/add-teams", nil)
	noBodyReq = mux.SetURLVars(noBodyReq, map[string]string{"id": "any"})

	throwReq := httptest.NewRequest(http.MethodPost, "/tournaments/:id/add-teams", nil)
	throwReq = mux.SetURLVars(throwReq, map[string]string{"id": "any"})

	throwReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	goodReq2 := httptest.NewRequest(http.MethodPost, "/tournaments/:id/add-teams", nil)
	goodReq2 = mux.SetURLVars(goodReq2, map[string]string{"id": "any"})
	goodReq2.Body = ioutil.NopCloser(bytes.NewReader(body))

	emptyIDReq := httptest.NewRequest(http.MethodPost, "/tournaments/:id/add-teams", nil)
	emptyIDReq = mux.SetURLVars(emptyIDReq, map[string]string{"id": ""})
	emptyIDReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	emptyBodyReq := httptest.NewRequest(http.MethodPost, "/tournaments/:id/add-teams", nil)
	emptyBodyReq = mux.SetURLVars(emptyBodyReq, map[string]string{"id": "any"})
	emptyBodyReq.Body = ioutil.NopCloser(bytes.NewReader(emptyBody))

	testCases := []struct {
		Name                       string
		Request                    *http.Request
		HandleUpdateTournamentFunc func(ctx context.Context, p tournament.Tournament) (*tournament.Tournament, errs.AppError)
		HandleGetTournamentFunc    func(ctx context.Context, id string) (*tournament.Tournament, errs.AppError)
		HandleGetTeamFunc          func(ctx context.Context, id string) (*team.Team, errs.AppError)
		ExpectedStatusCode         int
	}{
		{
			Name:                       "Should return 201 if successful",
			Request:                    goodReq,
			HandleUpdateTournamentFunc: mockUpdateTournamentFunc,
			HandleGetTournamentFunc:    mockGetTournamentFunc,
			HandleGetTeamFunc:          mockGetTeamFunc,
			ExpectedStatusCode:         200,
		}, {
			Name:                       "Should return 422 bad request",
			Request:                    noBodyReq,
			HandleUpdateTournamentFunc: mockUpdateTournamentFunc,
			HandleGetTournamentFunc:    mockGetTournamentFunc,
			HandleGetTeamFunc:          mockGetTeamFunc,
			ExpectedStatusCode:         422,
		}, {
			Name:                       "Should return 404 with empty id",
			Request:                    emptyIDReq,
			HandleUpdateTournamentFunc: mockUpdateTournamentFunc,
			HandleGetTournamentFunc:    mockGetTournamentFunc,
			HandleGetTeamFunc:          mockGetTeamFunc,
			ExpectedStatusCode:         404,
		}, {
			Name:                       "Should return 422 with body empty",
			Request:                    emptyBodyReq,
			HandleUpdateTournamentFunc: mockUpdateTournamentFunc,
			HandleGetTournamentFunc:    mockGetTournamentFunc,
			HandleGetTeamFunc:          mockGetTeamFunc,
			ExpectedStatusCode:         422,
		}, {
			Name:                       "Should return 500 throwing error on update tournament function",
			Request:                    throwReq,
			HandleUpdateTournamentFunc: mockUpdateTournamentThrowFunc,
			HandleGetTournamentFunc:    mockGetTournamentFunc,
			HandleGetTeamFunc:          mockGetTeamFunc,
			ExpectedStatusCode:         500,
		}, {
			Name:                       "Should return 500 throwing error on get tournament function",
			Request:                    goodReq2,
			HandleUpdateTournamentFunc: mockUpdateTournamentFunc,
			HandleGetTournamentFunc:    mockGetTournamentThrowFunc,
			HandleGetTeamFunc:          mockGetTeamFunc,
			ExpectedStatusCode:         500,
		}, {
			Name:                       "Should return 500 throwing error on get team function",
			Request:                    goodReq2,
			HandleUpdateTournamentFunc: mockUpdateTournamentFunc,
			HandleGetTournamentFunc:    mockGetTournamentFunc,
			HandleGetTeamFunc:          mockGetTeamThrowFunc,
			ExpectedStatusCode:         500,
		}, {
			Name:                       "Should return 404 throwing error on get tournament function returning nil",
			Request:                    goodReq2,
			HandleUpdateTournamentFunc: mockUpdateTournamentFunc,
			HandleGetTournamentFunc:    mockGetTournamentNilFunc,
			HandleGetTeamFunc:          mockGetTeamThrowFunc,
			ExpectedStatusCode:         404,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		repo.SetTournamentRepo(repo.MockTournamentRepo{
			UpdateFunc: tc.HandleUpdateTournamentFunc,
			GetFunc:    tc.HandleGetTournamentFunc,
		})
		defer repo.SetTournamentRepo(nil)

		repo.SetTeamRepo(repo.MockTeamRepo{
			GetFunc: tc.HandleGetTeamFunc,
		})
		defer repo.SetTeamRepo(nil)

		w := httptest.NewRecorder()

		HandleAddTeamsTournament(w, tc.Request)
		assert.Equal(t, tc.ExpectedStatusCode, w.Code)
	}
}
