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

func mockPostTeamFunc(ctx context.Context, t team.Team) errs.AppError {
	return nil
}

func mockPostTeamThrowFunc(ctx context.Context, t team.Team) errs.AppError {
	return errs.ErrRepoMockAction
}

func TestHandlePostTeam(t *testing.T) {
	body, err := json.Marshal(TeamEntityPayload{
		Name:      "Real Madrid Club",
		ShortCode: "RMC",
		Country:   "Spain",
		City:      "Madrid",
	})
	assert.Equal(t, nil, err)

	goodReq := httptest.NewRequest(http.MethodPost, "/teams", nil)
	goodReq = mux.SetURLVars(goodReq, map[string]string{})
	goodReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	noBodyReq := httptest.NewRequest(http.MethodPost, "/teams", nil)
	noBodyReq = mux.SetURLVars(noBodyReq, map[string]string{})

	throwReq := httptest.NewRequest(http.MethodPost, "/teams", nil)
	throwReq = mux.SetURLVars(throwReq, map[string]string{})

	throwReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	goodReq2 := httptest.NewRequest(http.MethodPost, "/teams", nil)
	goodReq2 = mux.SetURLVars(goodReq2, map[string]string{})
	goodReq2.Body = ioutil.NopCloser(bytes.NewReader(body))

	testCases := []struct {
		Name                  string
		Request               *http.Request
		HandlePostFunc        func(ctx context.Context, t team.Team) errs.AppError
		ConvertingPayloadFunc func(t TeamEntityPayload) (team.Team, errs.AppError)
		ExpectedStatusCode    int
	}{
		{
			Name:                  "Should return 201 if successful",
			Request:               goodReq,
			HandlePostFunc:        mockPostTeamFunc,
			ConvertingPayloadFunc: convertPayloadToTeamFunc,
			ExpectedStatusCode:    201,
		}, {
			Name:                  "Should return 422 bad request",
			Request:               noBodyReq,
			HandlePostFunc:        mockPostTeamFunc,
			ConvertingPayloadFunc: convertPayloadToTeamFunc,
			ExpectedStatusCode:    422,
		}, {
			Name:                  "Should return 422 throwing error on converting payload func",
			Request:               goodReq2,
			HandlePostFunc:        mockPostTeamFunc,
			ConvertingPayloadFunc: fakeConvertPayloadToTeamFunc,
			ExpectedStatusCode:    422,
		}, {
			Name:                  "Should return 500 throwing error on function",
			Request:               throwReq,
			HandlePostFunc:        mockPostTeamThrowFunc,
			ConvertingPayloadFunc: convertPayloadToTeamFunc,
			ExpectedStatusCode:    500,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		repo.SetTeamRepo(repo.MockTeamRepo{
			InsertFunc: tc.HandlePostFunc,
		})
		defer repo.SetTeamRepo(nil)

		convertPayloadToTeamFunc = tc.ConvertingPayloadFunc
		defer restoreConvertPayloadToTeamFunc(convertPayloadToTeamFunc)

		w := httptest.NewRecorder()

		HandlePostTeam(w, tc.Request)
		assert.Equal(t, tc.ExpectedStatusCode, w.Code)
	}
}

func TestConvertPayloadToTeam(t *testing.T) {
	inPayload := TeamEntityPayload{
		Name:      "Real Madrid Club",
		ShortCode: "RMC",
		Country:   "Spain",
		City:      "Madrid",
	}

	expectedTeam := team.Team{
		ID:        "",
		Name:      "Real Madrid Club",
		ShortCode: "RMC",
		Country:   "Spain",
		City:      "Madrid",
	}

	testCases := []struct {
		Name          string
		Payload       TeamEntityPayload
		ExpectedTeam  team.Team
		ExpectError   bool
		ExpectedError string
	}{
		{
			Name:         "Test Case: 1 - correct body, no error",
			Payload:      inPayload,
			ExpectedTeam: expectedTeam,
			ExpectError:  false,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		team, err := convertPayloadToTeam(tc.Payload)
		if tc.ExpectError {
			assert.NotNil(t, err)
			assert.Equal(t, err.Error(), tc.ExpectedError)
		} else {
			assert.Equal(t, tc.ExpectedTeam, team)
		}
	}
}

func TestDecodeTeamRequest(t *testing.T) {
	body, err := json.Marshal(TeamEntityPayload{
		Name:      "Real Madrid Club",
		ShortCode: "RMC",
		Country:   "Spain",
		City:      "Madrid",
	})
	assert.Equal(t, nil, err)

	goodReq := httptest.NewRequest(http.MethodPost, "/teams", nil)
	goodReq = mux.SetURLVars(goodReq, map[string]string{})

	goodReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	noBodyReq := httptest.NewRequest(http.MethodPost, "/teams", nil)
	noBodyReq = mux.SetURLVars(noBodyReq, map[string]string{})

	testCases := []struct {
		Name          string
		Request       *http.Request
		Payload       *TeamEntityPayload
		ExpectedError bool
	}{
		{
			Name:    "Test Case: 1 - correct body, no error",
			Request: goodReq, Payload: &TeamEntityPayload{
				Name:      "Real Madrid Club",
				ShortCode: "RMC",
				Country:   "Spain",
				City:      "Madrid",
			}, ExpectedError: false,
		},
		{Name: "Test Case: 2 - no body, error found", Request: noBodyReq, Payload: nil, ExpectedError: true},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		decodedPayload, err := decodeTeamRequest(tc.Request)
		if tc.ExpectedError {
			assert.NotNil(t, err)
		} else {
			assert.Equal(t, *tc.Payload, decodedPayload)
		}
	}
}
