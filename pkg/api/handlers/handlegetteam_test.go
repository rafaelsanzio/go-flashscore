package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
	"github.com/rafaelsanzio/go-flashscore/pkg/prototype"
	"github.com/rafaelsanzio/go-flashscore/pkg/team"
)

func mockGetTeamFunc(ctx context.Context, id string) (*team.Team, errs.AppError) {
	teamMock := prototype.PrototypeTeam()
	return &teamMock, nil
}

func mockGetTeamThrowFunc(ctx context.Context, id string) (*team.Team, errs.AppError) {
	return nil, errs.ErrRepoMockAction
}

func mockGetTeamNilFunc(ctx context.Context, id string) (*team.Team, errs.AppError) {
	return nil, nil
}

func TestHandleGetTeam(t *testing.T) {
	testCases := []struct {
		Name               string
		ID                 string
		HandleGetTeamFunc  func(ctx context.Context, id string) (*team.Team, errs.AppError)
		MarshalFunc        func(v interface{}) ([]byte, error)
		WriteFunc          func(http.ResponseWriter, []byte) (int, error)
		ExpectedStatusCode int
	}{
		{
			Name:               "Success handle get team",
			ID:                 "1",
			HandleGetTeamFunc:  mockGetTeamFunc,
			MarshalFunc:        jsonMarshal,
			WriteFunc:          write,
			ExpectedStatusCode: 200,
		}, {
			Name:               "Not Found handle get team",
			ID:                 "",
			HandleGetTeamFunc:  mockGetTeamFunc,
			MarshalFunc:        jsonMarshal,
			WriteFunc:          write,
			ExpectedStatusCode: 404,
		}, {
			Name:               "Getting error on team repo",
			ID:                 "1",
			HandleGetTeamFunc:  mockGetTeamThrowFunc,
			MarshalFunc:        jsonMarshal,
			WriteFunc:          write,
			ExpectedStatusCode: 500,
		}, {
			Name:               "Getting error on marshal function",
			ID:                 "1",
			HandleGetTeamFunc:  mockGetTeamFunc,
			MarshalFunc:        fakeMarshal,
			WriteFunc:          write,
			ExpectedStatusCode: 500,
		}, {
			Name:               "Getting error on write function",
			ID:                 "1",
			HandleGetTeamFunc:  mockGetTeamFunc,
			MarshalFunc:        jsonMarshal,
			WriteFunc:          fakeWrite,
			ExpectedStatusCode: 500,
		}, {
			Name:               "Getting error on get func returning nil",
			ID:                 "1",
			HandleGetTeamFunc:  mockGetTeamNilFunc,
			MarshalFunc:        jsonMarshal,
			WriteFunc:          write,
			ExpectedStatusCode: 404,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		repo.SetTeamRepo(repo.MockTeamRepo{
			GetFunc: tc.HandleGetTeamFunc,
		})
		defer repo.SetTeamRepo(nil)

		jsonMarshal = tc.MarshalFunc
		defer restoreMarshal(jsonMarshal)

		write = tc.WriteFunc
		defer restoreWrite(write)

		req, err := http.NewRequest(http.MethodGet, "/teams/:id", nil)
		req = mux.SetURLVars(req, map[string]string{"id": tc.ID})
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		HandleGetTeam(res, req)

		assert.Equal(t, tc.ExpectedStatusCode, res.Code)
		t.Logf("Response Body: %v", res.Body)

		if res.Code == http.StatusOK {
			team := team.Team{}
			err = json.Unmarshal(res.Body.Bytes(), &team)
			assert.NoError(t, err)
		}
	}
}
