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
	"github.com/rafaelsanzio/go-flashscore/pkg/player"
)

func mockDeletePlayerFunc(ctx context.Context, id string) errs.AppError {
	return nil
}

func mockDeletePlayerThrowFunc(ctx context.Context, id string) errs.AppError {
	return errs.ErrRepoMockAction
}

func TestHandleDeletePlayer(t *testing.T) {
	testCases := []struct {
		Name                   string
		ID                     string
		HandleDeletePlayerFunc func(ctx context.Context, id string) errs.AppError
		HandleGetPlayerFunc    func(ctx context.Context, id string) (*player.Player, errs.AppError)
		ExpectedStatusCode     int
	}{
		{
			Name:                   "Success handle delete player",
			ID:                     "1",
			HandleDeletePlayerFunc: mockDeletePlayerFunc,
			HandleGetPlayerFunc:    mockGetPlayerFunc,
			ExpectedStatusCode:     204,
		}, {
			Name:                   "Not Found handle delete player",
			ID:                     "",
			HandleDeletePlayerFunc: mockDeletePlayerFunc,
			HandleGetPlayerFunc:    mockGetPlayerFunc,
			ExpectedStatusCode:     404,
		}, {
			Name:                   "Throwing error on delete function",
			ID:                     "1",
			HandleDeletePlayerFunc: mockDeletePlayerThrowFunc,
			HandleGetPlayerFunc:    mockGetPlayerFunc,
			ExpectedStatusCode:     500,
		}, {
			Name:                   "Throwing error on get function",
			ID:                     "1",
			HandleDeletePlayerFunc: mockDeletePlayerFunc,
			HandleGetPlayerFunc:    mockGetPlayerThrowFunc,
			ExpectedStatusCode:     500,
		}, {
			Name:                   "Throwing error on get function returning nil",
			ID:                     "1",
			HandleDeletePlayerFunc: mockDeletePlayerFunc,
			HandleGetPlayerFunc:    mockGetPlayerNilFunc,
			ExpectedStatusCode:     404,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		repo.SetPlayerRepo(repo.MockPlayerRepo{
			DeleteFunc: tc.HandleDeletePlayerFunc,
			GetFunc:    tc.HandleGetPlayerFunc,
		})
		defer repo.SetPlayerRepo(nil)

		req, err := http.NewRequest(http.MethodDelete, "/player/:id", nil)
		req = mux.SetURLVars(req, map[string]string{"id": tc.ID})
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		HandleDeletePlayer(res, req)

		assert.Equal(t, tc.ExpectedStatusCode, res.Code)

		if res.Code == http.StatusNoContent {
			assert.NoError(t, err)
		}
	}
}
