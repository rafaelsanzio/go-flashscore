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
	"github.com/rafaelsanzio/go-flashscore/pkg/player"
	"github.com/rafaelsanzio/go-flashscore/pkg/team"
)

func mockUpdatePlayerFunc(ctx context.Context, p player.Player) (*player.Player, errs.AppError) {
	return &p, nil
}

func mockUpdatePlayerThrowFunc(ctx context.Context, p player.Player) (*player.Player, errs.AppError) {
	return nil, errs.ErrRepoMockAction
}

func TestHandleUpdatePlayer(t *testing.T) {
	body, err := json.Marshal(PlayerEntityPayload{
		Name:         "Cristiano Ronaldo",
		Team:         "any_team_id",
		Country:      "Portugal",
		BirthdayDate: "1990-01-01",
	})
	assert.Equal(t, nil, err)

	goodReq := httptest.NewRequest(http.MethodPut, "/players/:id", nil)
	goodReq = mux.SetURLVars(goodReq, map[string]string{"id": "any_id"})
	goodReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	noBodyReq := httptest.NewRequest(http.MethodPut, "/players/:id", nil)
	noBodyReq = mux.SetURLVars(noBodyReq, map[string]string{"id": "any_id"})

	throwReq := httptest.NewRequest(http.MethodPut, "/players/:id", nil)
	throwReq = mux.SetURLVars(throwReq, map[string]string{"id": "any_id"})
	throwReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	goodReq2 := httptest.NewRequest(http.MethodPut, "/players/:id", nil)
	goodReq2 = mux.SetURLVars(goodReq2, map[string]string{"id": "any_id"})
	goodReq2.Body = ioutil.NopCloser(bytes.NewReader(body))

	missParamReq := httptest.NewRequest(http.MethodPut, "/players/:id", nil)
	missParamReq = mux.SetURLVars(missParamReq, map[string]string{})
	missParamReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	testCases := []struct {
		Name                  string
		Request               *http.Request
		HandleUpdateFunction  func(ctx context.Context, t player.Player) (*player.Player, errs.AppError)
		HandleGetTeamFunc     func(ctx context.Context, id string) (*team.Team, errs.AppError)
		ConvertingPayloadFunc func(ctx context.Context, p PlayerEntityPayload) (*player.Player, errs.AppError)
		ExpectedStatusCode    int
	}{
		{
			Name:                  "Should return 200 if successful",
			Request:               goodReq,
			HandleUpdateFunction:  mockUpdatePlayerFunc,
			HandleGetTeamFunc:     mockGetTeamFunc,
			ConvertingPayloadFunc: convertPayloadToPlayerFunc,
			ExpectedStatusCode:    200,
		}, {
			Name:                  "Throwing error on function",
			Request:               throwReq,
			HandleUpdateFunction:  mockUpdatePlayerThrowFunc,
			HandleGetTeamFunc:     mockGetTeamFunc,
			ConvertingPayloadFunc: convertPayloadToPlayerFunc,
			ExpectedStatusCode:    500,
		}, {
			Name:                  "Should return 422 bad request",
			Request:               noBodyReq,
			HandleUpdateFunction:  mockUpdatePlayerFunc,
			HandleGetTeamFunc:     mockGetTeamFunc,
			ConvertingPayloadFunc: convertPayloadToPlayerFunc,
			ExpectedStatusCode:    422,
		}, {
			Name:                  "Should return 500 throwing error on convertPayloadToPlayer function",
			Request:               goodReq2,
			HandleUpdateFunction:  mockUpdatePlayerFunc,
			HandleGetTeamFunc:     mockGetTeamFunc,
			ConvertingPayloadFunc: fakeConvertPayloadToPlayerFunc,
			ExpectedStatusCode:    422,
		}, {
			Name:                  "Should return 404 is missing param",
			Request:               missParamReq,
			HandleUpdateFunction:  mockUpdatePlayerFunc,
			HandleGetTeamFunc:     mockGetTeamFunc,
			ConvertingPayloadFunc: convertPayloadToPlayerFunc,
			ExpectedStatusCode:    404,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		repo.SetPlayerRepo(repo.MockPlayerRepo{
			UpdateFunc: tc.HandleUpdateFunction,
		})
		defer repo.SetPlayerRepo(nil)

		repo.SetTeamRepo(repo.MockTeamRepo{
			GetFunc: tc.HandleGetTeamFunc,
		})
		defer repo.SetTeamRepo(nil)

		convertPayloadToPlayerFunc = tc.ConvertingPayloadFunc
		defer restoreConvertPayloadToPlayerFunc(convertPayloadToPlayerFunc)

		w := httptest.NewRecorder()

		HandleUpdatePlayer(w, tc.Request)

		assert.Equal(t, tc.ExpectedStatusCode, w.Code)
	}
}
