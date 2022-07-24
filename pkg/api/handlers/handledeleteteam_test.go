package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
	"github.com/rafaelsanzio/go-flashscore/pkg/team"
)

func mockDeleteTeamFunc(ctx context.Context, id string) errs.AppError {
	return nil
}

func mockDeleteTeamThrowFunc(ctx context.Context, id string) errs.AppError {
	return errs.ErrRepoMockAction
}

func TestHandleDeleteTeam(t *testing.T) {
	testCases := []struct {
		Name                 string
		ID                   string
		HandleDeleteTeamFunc func(ctx context.Context, id string) errs.AppError
		HandleGetTeamFunc    func(ctx context.Context, id string) (*team.Team, errs.AppError)
		ExpectedStatusCode   int
	}{
		{
			Name:                 "Success handle delete team",
			ID:                   "1",
			HandleDeleteTeamFunc: mockDeleteTeamFunc,
			HandleGetTeamFunc:    mockGetTeamFunc,
			ExpectedStatusCode:   204,
		}, {
			Name:                 "Not Found handle delete team",
			ID:                   "",
			HandleDeleteTeamFunc: mockDeleteTeamFunc,
			HandleGetTeamFunc:    mockGetTeamFunc,
			ExpectedStatusCode:   404,
		}, {
			Name:                 "Throwing error on delete function",
			ID:                   "1",
			HandleDeleteTeamFunc: mockDeleteTeamThrowFunc,
			HandleGetTeamFunc:    mockGetTeamFunc,
			ExpectedStatusCode:   500,
		}, {
			Name:                 "Throwing error on get function",
			ID:                   "1",
			HandleDeleteTeamFunc: mockDeleteTeamFunc,
			HandleGetTeamFunc:    mockGetTeamThrowFunc,
			ExpectedStatusCode:   500,
		}, {
			Name:                 "Throwing error on get function returning nil",
			ID:                   "1",
			HandleDeleteTeamFunc: mockDeleteTeamFunc,
			HandleGetTeamFunc:    mockGetTeamNilFunc,
			ExpectedStatusCode:   404,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		repo.SetTeamRepo(repo.MockTeamRepo{
			DeleteFunc: tc.HandleDeleteTeamFunc,
			GetFunc:    tc.HandleGetTeamFunc,
		})
		defer repo.SetTeamRepo(nil)

		req, err := http.NewRequest(http.MethodDelete, "/teams/:id", nil)
		req = mux.SetURLVars(req, map[string]string{"id": tc.ID})
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		HandleDeleteTeam(res, req)

		assert.Equal(t, tc.ExpectedStatusCode, res.Code)

		if res.Code == 204 {
			assert.NoError(t, err)
		}
	}
}
