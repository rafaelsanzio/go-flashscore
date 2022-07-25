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
)

func mockUpdateTeamFunc(ctx context.Context, t team.Team) (*team.Team, errs.AppError) {
	return &t, nil
}

func mockUpdateTeamThrowFunc(ctx context.Context, t team.Team) (*team.Team, errs.AppError) {
	return nil, errs.ErrRepoMockAction
}

func TestHandleUpdateTeam(t *testing.T) {
	body, err := json.Marshal(TeamEntityPayload{
		Name:      "Real Madrid Club",
		ShortCode: "RMC",
		Country:   "Spain",
		City:      "Madrid",
	})
	assert.Equal(t, nil, err)

	goodReq := httptest.NewRequest(http.MethodPut, "/teams/:id", nil)
	goodReq = mux.SetURLVars(goodReq, map[string]string{"id": "any_id"})
	goodReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	noBodyReq := httptest.NewRequest(http.MethodPut, "/teams/:id", nil)
	noBodyReq = mux.SetURLVars(noBodyReq, map[string]string{"id": "any_id"})

	throwReq := httptest.NewRequest(http.MethodPut, "/teams/:id", nil)
	throwReq = mux.SetURLVars(throwReq, map[string]string{"id": "any_id"})
	throwReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	goodReq2 := httptest.NewRequest(http.MethodPut, "/teams/:id", nil)
	goodReq2 = mux.SetURLVars(goodReq2, map[string]string{"id": "any_id"})
	goodReq2.Body = ioutil.NopCloser(bytes.NewReader(body))

	missParamReq := httptest.NewRequest(http.MethodPut, "/teams/:id", nil)
	missParamReq = mux.SetURLVars(missParamReq, map[string]string{})
	missParamReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	testCases := []struct {
		Name                  string
		Request               *http.Request
		HandleUpdateFunction  func(ctx context.Context, t team.Team) (*team.Team, errs.AppError)
		ConvertingPayloadFunc func(t TeamEntityPayload) (team.Team, errs.AppError)
		ExpectedStatusCode    int
	}{
		{
			Name:                  "Should return 200 if successful",
			Request:               goodReq,
			HandleUpdateFunction:  mockUpdateTeamFunc,
			ConvertingPayloadFunc: convertPayloadToTeamFunc,
			ExpectedStatusCode:    200,
		}, {
			Name:                  "Throwing error on function",
			Request:               throwReq,
			HandleUpdateFunction:  mockUpdateTeamThrowFunc,
			ConvertingPayloadFunc: convertPayloadToTeamFunc,
			ExpectedStatusCode:    500,
		}, {
			Name:                  "Should return 422 bad request",
			Request:               noBodyReq,
			HandleUpdateFunction:  mockUpdateTeamFunc,
			ConvertingPayloadFunc: convertPayloadToTeamFunc,
			ExpectedStatusCode:    422,
		}, {
			Name:                  "Should return 500 throwing error on convertPayloadToTeam function",
			Request:               goodReq2,
			HandleUpdateFunction:  mockUpdateTeamFunc,
			ConvertingPayloadFunc: fakeConvertPayloadToTeamFunc,
			ExpectedStatusCode:    422,
		}, {
			Name:                  "Should return 404 is missing param",
			Request:               missParamReq,
			HandleUpdateFunction:  mockUpdateTeamFunc,
			ConvertingPayloadFunc: convertPayloadToTeamFunc,
			ExpectedStatusCode:    404,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		repo.SetTeamRepo(repo.MockTeamRepo{
			UpdateFunc: tc.HandleUpdateFunction,
		})
		defer repo.SetTeamRepo(nil)

		convertPayloadToTeamFunc = tc.ConvertingPayloadFunc
		defer restoreConvertPayloadToTeamFunc(convertPayloadToTeamFunc)

		w := httptest.NewRecorder()

		HandleUpdateTeam(w, tc.Request)

		assert.Equal(t, tc.ExpectedStatusCode, w.Code)
	}
}
