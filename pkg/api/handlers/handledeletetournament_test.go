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
	"github.com/rafaelsanzio/go-flashscore/pkg/tournament"
)

func mockDeleteTournamentFunc(ctx context.Context, id string) errs.AppError {
	return nil
}

func mockDeleteTournamentThrowFunc(ctx context.Context, id string) errs.AppError {
	return errs.ErrRepoMockAction
}

func TestHandleDeleteTournament(t *testing.T) {
	testCases := []struct {
		Name                       string
		ID                         string
		HandleDeleteTournemantFunc func(ctx context.Context, id string) errs.AppError
		HandleGetTournemantFunc    func(ctx context.Context, id string) (*tournament.Tournament, errs.AppError)
		ExpectedStatusCode         int
	}{
		{
			Name:                       "Success handle delete tournament",
			ID:                         "1",
			HandleDeleteTournemantFunc: mockDeleteTournamentFunc,
			HandleGetTournemantFunc:    mockGetTournamentFunc,
			ExpectedStatusCode:         204,
		}, {
			Name:                       "Not Found handle delete tournament",
			ID:                         "",
			HandleDeleteTournemantFunc: mockDeleteTournamentFunc,
			HandleGetTournemantFunc:    mockGetTournamentFunc,
			ExpectedStatusCode:         404,
		}, {
			Name:                       "Throwing error on delete function",
			ID:                         "1",
			HandleDeleteTournemantFunc: mockDeleteTournamentThrowFunc,
			HandleGetTournemantFunc:    mockGetTournamentFunc,
			ExpectedStatusCode:         500,
		}, {
			Name:                       "Throwing error on get function",
			ID:                         "1",
			HandleDeleteTournemantFunc: mockDeleteTournamentFunc,
			HandleGetTournemantFunc:    mockGetTournamentThrowFunc,
			ExpectedStatusCode:         500,
		}, {
			Name:                       "Throwing error on get function returning nil",
			ID:                         "1",
			HandleDeleteTournemantFunc: mockDeleteTournamentFunc,
			HandleGetTournemantFunc:    mockGetTournamentNilFunc,
			ExpectedStatusCode:         404,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		repo.SetTournamentRepo(repo.MockTournamentRepo{
			DeleteFunc: tc.HandleDeleteTournemantFunc,
			GetFunc:    tc.HandleGetTournemantFunc,
		})
		defer repo.SetTournamentRepo(nil)

		req, err := http.NewRequest(http.MethodDelete, "/tournaments/:id", nil)
		req = mux.SetURLVars(req, map[string]string{"id": tc.ID})
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		HandleDeleteTournament(res, req)

		assert.Equal(t, tc.ExpectedStatusCode, res.Code)

		if res.Code == http.StatusNoContent {
			assert.NoError(t, err)
		}
	}
}
