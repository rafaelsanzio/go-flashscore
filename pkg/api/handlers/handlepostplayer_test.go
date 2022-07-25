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
	"github.com/rafaelsanzio/go-flashscore/pkg/prototype"
	"github.com/rafaelsanzio/go-flashscore/pkg/team"
)

func mockPostPlayerFunc(ctx context.Context, p player.Player) errs.AppError {
	return nil
}

func mockPostPlayerThrowFunc(ctx context.Context, p player.Player) errs.AppError {
	return errs.ErrRepoMockAction
}

func TestHandlePostPlayer(t *testing.T) {
	body, err := json.Marshal(PlayerEntityPayload{
		Name:         "Cristiano Ronaldo",
		Team:         "any_team_id",
		Country:      "Portugal",
		BirthdayDate: "1990-01-01",
	})
	assert.Equal(t, nil, err)

	goodReq := httptest.NewRequest(http.MethodPost, "/players", nil)
	goodReq = mux.SetURLVars(goodReq, map[string]string{})
	goodReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	noBodyReq := httptest.NewRequest(http.MethodPost, "/players", nil)
	noBodyReq = mux.SetURLVars(noBodyReq, map[string]string{})

	throwReq := httptest.NewRequest(http.MethodPost, "/players", nil)
	throwReq = mux.SetURLVars(throwReq, map[string]string{})

	throwReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	goodReq2 := httptest.NewRequest(http.MethodPost, "/players", nil)
	goodReq2 = mux.SetURLVars(goodReq2, map[string]string{})
	goodReq2.Body = ioutil.NopCloser(bytes.NewReader(body))

	testCases := []struct {
		Name                  string
		Request               *http.Request
		HandlePostFunc        func(ctx context.Context, p player.Player) errs.AppError
		HandleGetTeamFunc     func(ctx context.Context, id string) (*team.Team, errs.AppError)
		ConvertingPayloadFunc func(ctx context.Context, p PlayerEntityPayload) (*player.Player, errs.AppError)
		ExpectedStatusCode    int
	}{
		{
			Name:                  "Should return 201 if successful",
			Request:               goodReq,
			HandlePostFunc:        mockPostPlayerFunc,
			HandleGetTeamFunc:     mockGetTeamFunc,
			ConvertingPayloadFunc: convertPayloadToPlayerFunc,
			ExpectedStatusCode:    201,
		}, {
			Name:                  "Should return 422 bad request",
			Request:               noBodyReq,
			HandlePostFunc:        mockPostPlayerFunc,
			HandleGetTeamFunc:     mockGetTeamFunc,
			ConvertingPayloadFunc: convertPayloadToPlayerFunc,
			ExpectedStatusCode:    422,
		}, {
			Name:                  "Should return 422 throwing error on converting payload func",
			Request:               goodReq2,
			HandlePostFunc:        mockPostPlayerFunc,
			HandleGetTeamFunc:     mockGetTeamFunc,
			ConvertingPayloadFunc: fakeConvertPayloadToPlayerFunc,
			ExpectedStatusCode:    422,
		}, {
			Name:                  "Should return 500 throwing error on function",
			Request:               throwReq,
			HandlePostFunc:        mockPostPlayerThrowFunc,
			HandleGetTeamFunc:     mockGetTeamFunc,
			ConvertingPayloadFunc: convertPayloadToPlayerFunc,
			ExpectedStatusCode:    500,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		repo.SetPlayerRepo(repo.MockPlayerRepo{
			InsertFunc: tc.HandlePostFunc,
		})
		defer repo.SetPlayerRepo(nil)

		repo.SetTeamRepo(repo.MockTeamRepo{
			GetFunc: tc.HandleGetTeamFunc,
		})
		defer repo.SetTeamRepo(nil)

		convertPayloadToPlayerFunc = tc.ConvertingPayloadFunc
		defer restoreConvertPayloadToPlayerFunc(convertPayloadToPlayerFunc)

		w := httptest.NewRecorder()

		HandlePostPlayer(w, tc.Request)
		assert.Equal(t, tc.ExpectedStatusCode, w.Code)
	}
}

func TestConvertPayloadToPlayer(t *testing.T) {
	inPayload := PlayerEntityPayload{
		Name:         "Cristiano Ronaldo",
		Team:         "any_team_id",
		Country:      "Portugal",
		BirthdayDate: "1990-01-01",
	}

	expectedTeam := player.Player{
		ID:           "",
		Name:         "Cristiano Ronaldo",
		Team:         prototype.PrototypeTeam(),
		Country:      "Portugal",
		BirthdayDate: "1990-01-01",
	}

	testCases := []struct {
		Name          string
		Payload       PlayerEntityPayload
		HandleGetFunc func(ctx context.Context, id string) (*team.Team, errs.AppError)
		ExpectedTeam  player.Player
		ExpectError   bool
	}{
		{
			Name:          "Test Case: 1 - correct body, no error",
			Payload:       inPayload,
			HandleGetFunc: mockGetTeamFunc,
			ExpectedTeam:  expectedTeam,
			ExpectError:   false,
		}, {
			Name:          "Test Case: 2 - throwing error on get function",
			Payload:       inPayload,
			HandleGetFunc: mockGetTeamThrowFunc,
			ExpectedTeam:  expectedTeam,
			ExpectError:   true,
		}, {
			Name:          "Test Case: 3 - throwing error on get team nil function",
			Payload:       inPayload,
			HandleGetFunc: mockGetTeamNilFunc,
			ExpectedTeam:  expectedTeam,
			ExpectError:   true,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		repo.SetTeamRepo(repo.MockTeamRepo{
			GetFunc: tc.HandleGetFunc,
		})
		defer repo.SetTeamRepo(nil)

		player, err := convertPayloadToPlayerFunc(context.Background(), tc.Payload)
		if tc.ExpectError {
			assert.NotNil(t, err)
		} else {
			assert.Equal(t, tc.ExpectedTeam, *player)
		}
	}
}
