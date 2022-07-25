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

func mockUpdateTournamentFunc(ctx context.Context, t tournament.Tournament) (*tournament.Tournament, errs.AppError) {
	return &t, nil
}

func mockUpdateTournamentThrowFunc(ctx context.Context, t tournament.Tournament) (*tournament.Tournament, errs.AppError) {
	return nil, errs.ErrRepoMockAction
}

func TestHandleUpdateTournament(t *testing.T) {
	body, err := json.Marshal(TournamentEntityPayload{
		Name:  "Any Tournament Name",
		Teams: []string{"any_team_id", "any_team_id_2"},
	})
	assert.Equal(t, nil, err)

	goodReq := httptest.NewRequest(http.MethodPut, "/tournaments/:id", nil)
	goodReq = mux.SetURLVars(goodReq, map[string]string{"id": "any_id"})
	goodReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	noBodyReq := httptest.NewRequest(http.MethodPut, "/tournaments/:id", nil)
	noBodyReq = mux.SetURLVars(noBodyReq, map[string]string{"id": "any_id"})

	throwReq := httptest.NewRequest(http.MethodPut, "/tournaments/:id", nil)
	throwReq = mux.SetURLVars(throwReq, map[string]string{"id": "any_id"})
	throwReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	goodReq2 := httptest.NewRequest(http.MethodPut, "/tournaments/:id", nil)
	goodReq2 = mux.SetURLVars(goodReq2, map[string]string{"id": "any_id"})
	goodReq2.Body = ioutil.NopCloser(bytes.NewReader(body))

	missParamReq := httptest.NewRequest(http.MethodPut, "/tournaments/:id", nil)
	missParamReq = mux.SetURLVars(missParamReq, map[string]string{})
	missParamReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	testCases := []struct {
		Name                           string
		Request                        *http.Request
		HandleUpdateTournamentFunction func(ctx context.Context, t tournament.Tournament) (*tournament.Tournament, errs.AppError)
		HandleGetTeamFunc              func(ctx context.Context, id string) (*team.Team, errs.AppError)
		ConvertingPayloadFunc          func(ctx context.Context, p TournamentEntityPayload) (*tournament.Tournament, errs.AppError)
		ExpectedStatusCode             int
	}{
		{
			Name:                           "Should return 200 if successful",
			Request:                        goodReq,
			HandleUpdateTournamentFunction: mockUpdateTournamentFunc,
			HandleGetTeamFunc:              mockGetTeamFunc,
			ConvertingPayloadFunc:          convertPayloadToTournamentFunc,
			ExpectedStatusCode:             200,
		}, {
			Name:                           "Throwing error on function",
			Request:                        throwReq,
			HandleUpdateTournamentFunction: mockUpdateTournamentThrowFunc,
			HandleGetTeamFunc:              mockGetTeamFunc,
			ConvertingPayloadFunc:          convertPayloadToTournamentFunc,
			ExpectedStatusCode:             500,
		}, {
			Name:                           "Should return 422 bad request",
			Request:                        noBodyReq,
			HandleUpdateTournamentFunction: mockUpdateTournamentFunc,
			HandleGetTeamFunc:              mockGetTeamFunc,
			ConvertingPayloadFunc:          convertPayloadToTournamentFunc,
			ExpectedStatusCode:             422,
		}, {
			Name:                           "Should return 500 throwing error on convertPayloadToPlayer function",
			Request:                        goodReq2,
			HandleUpdateTournamentFunction: mockUpdateTournamentFunc,
			HandleGetTeamFunc:              mockGetTeamFunc,
			ConvertingPayloadFunc:          fakeConvertPayloadToTournament,
			ExpectedStatusCode:             422,
		}, {
			Name:                           "Should return 404 is missing param",
			Request:                        missParamReq,
			HandleUpdateTournamentFunction: mockUpdateTournamentFunc,
			HandleGetTeamFunc:              mockGetTeamFunc,
			ConvertingPayloadFunc:          convertPayloadToTournamentFunc,
			ExpectedStatusCode:             404,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		repo.SetTournamentRepo(repo.MockTournamentRepo{
			UpdateFunc: tc.HandleUpdateTournamentFunction,
		})
		defer repo.SetTournamentRepo(nil)

		repo.SetTeamRepo(repo.MockTeamRepo{
			GetFunc: tc.HandleGetTeamFunc,
		})
		defer repo.SetTeamRepo(nil)

		convertPayloadToTournamentFunc = tc.ConvertingPayloadFunc
		defer restoreConvertPayloadToTournament(convertPayloadToTournamentFunc)

		w := httptest.NewRecorder()

		HandleUpdateTournament(w, tc.Request)

		assert.Equal(t, tc.ExpectedStatusCode, w.Code)
	}
}
