package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
	"github.com/rafaelsanzio/go-flashscore/pkg/prototype"
	"github.com/rafaelsanzio/go-flashscore/pkg/team"
)

func mockListTeamFunc(ctx context.Context) ([]team.Team, errs.AppError) {
	teamMock := prototype.PrototypeTeam()

	teamMock2 := prototype.PrototypeTeam()
	teamMock2.Name = "Real Madrid B"

	teamMockList := []team.Team{teamMock, teamMock2}

	return teamMockList, nil
}

func mockListTeamThrowFunc(ctx context.Context) ([]team.Team, errs.AppError) {
	return nil, errs.ErrRepoMockAction
}

func TestHandleListTeam(t *testing.T) {
	testCases := []struct {
		Name               string
		HandleListTeamFunc func(ctx context.Context) ([]team.Team, errs.AppError)
		MarshalFunc        func(v interface{}) ([]byte, error)
		WriteFunc          func(http.ResponseWriter, []byte) (int, error)
		ExpectedStatusCode int
	}{
		{
			Name:               "Success handle list teams",
			HandleListTeamFunc: mockListTeamFunc,
			MarshalFunc:        jsonMarshal,
			WriteFunc:          write,
			ExpectedStatusCode: 200,
		}, {
			Name:               "Throwing handle list teams",
			HandleListTeamFunc: mockListTeamThrowFunc,
			MarshalFunc:        jsonMarshal,
			WriteFunc:          write,
			ExpectedStatusCode: 500,
		}, {
			Name:               "Throwing error on marshal function",
			HandleListTeamFunc: mockListTeamFunc,
			MarshalFunc:        fakeMarshal,
			WriteFunc:          write,
			ExpectedStatusCode: 500,
		}, {
			Name:               "Throwing error on write function",
			HandleListTeamFunc: mockListTeamFunc,
			MarshalFunc:        jsonMarshal,
			WriteFunc:          fakeWrite,
			ExpectedStatusCode: 500,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		repo.SetTeamRepo(repo.MockTeamRepo{
			ListFunc: tc.HandleListTeamFunc,
		})
		defer repo.SetTeamRepo(nil)

		jsonMarshal = tc.MarshalFunc
		defer restoreMarshal(jsonMarshal)

		write = tc.WriteFunc
		defer restoreWrite(write)

		req, err := http.NewRequest(http.MethodGet, "/teams", nil)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		HandleListTeam(res, req)

		assert.Equal(t, tc.ExpectedStatusCode, res.Code)
		t.Logf("Response Body: %v", res.Body)

		if res.Code == http.StatusOK {
			team := []team.Team{}
			err = json.Unmarshal(res.Body.Bytes(), &team)
			assert.NoError(t, err)

			assert.Equal(t, 2, len(team))
		}
	}
}
