package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/match"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
	"github.com/rafaelsanzio/go-flashscore/pkg/tournament"
)

func mockDeleteMatchFunc(ctx context.Context, id string) errs.AppError {
	return nil
}

func mockDeleteMatchThrowFunc(ctx context.Context, id string) errs.AppError {
	return errs.ErrRepoMockAction
}

func TestHandleDeleteMatch(t *testing.T) {
	testCases := []struct {
		Name                             string
		ID                               string
		MatchID                          string
		HandleDeleteMatchFunc            func(ctx context.Context, id string) errs.AppError
		HandleFindMatchForTournamentFunc func(ctx context.Context, id, tournamentID string) (*match.Match, errs.AppError)
		HandleGetTournemantFunc          func(ctx context.Context, id string) (*tournament.Tournament, errs.AppError)
		ExpectedStatusCode               int
	}{
		{
			Name:                             "Success handle delete match",
			ID:                               "1",
			MatchID:                          "1",
			HandleDeleteMatchFunc:            mockDeleteMatchFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			HandleGetTournemantFunc:          mockGetTournamentFunc,
			ExpectedStatusCode:               204,
		}, {
			Name:                             "Not Found id param handle delete match",
			ID:                               "",
			MatchID:                          "1",
			HandleDeleteMatchFunc:            mockDeleteMatchFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			HandleGetTournemantFunc:          mockGetTournamentFunc,
			ExpectedStatusCode:               404,
		}, {
			Name:                             "Not Found match_id param handle delete match",
			ID:                               "1",
			MatchID:                          "",
			HandleDeleteMatchFunc:            mockDeleteMatchFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			HandleGetTournemantFunc:          mockGetTournamentFunc,
			ExpectedStatusCode:               404,
		}, {
			Name:                             "Throwing error on delete function",
			ID:                               "1",
			MatchID:                          "1",
			HandleDeleteMatchFunc:            mockDeleteMatchThrowFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			HandleGetTournemantFunc:          mockGetTournamentFunc,
			ExpectedStatusCode:               500,
		}, {
			Name:                             "Throwing error on get function",
			ID:                               "1",
			MatchID:                          "1",
			HandleDeleteMatchFunc:            mockDeleteMatchFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentThrowFunc,
			HandleGetTournemantFunc:          mockGetTournamentFunc,
			ExpectedStatusCode:               500,
		}, {
			Name:                             "Throwing error on get function returning nil",
			ID:                               "1",
			MatchID:                          "1",
			HandleDeleteMatchFunc:            mockDeleteMatchFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentNilFunc,
			HandleGetTournemantFunc:          mockGetTournamentFunc,
			ExpectedStatusCode:               404,
		}, {
			Name:                             "Throwing error on get tournament function",
			ID:                               "1",
			MatchID:                          "1",
			HandleDeleteMatchFunc:            mockDeleteMatchThrowFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentNilFunc,
			HandleGetTournemantFunc:          mockGetTournamentThrowFunc,
			ExpectedStatusCode:               500,
		}, {
			Name:                             "Throwing error on get tournament function returning nil",
			ID:                               "1",
			MatchID:                          "1",
			HandleDeleteMatchFunc:            mockDeleteMatchThrowFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentNilFunc,
			HandleGetTournemantFunc:          mockGetTournamentNilFunc,
			ExpectedStatusCode:               404,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		repo.SetMatchRepo(repo.MockMatchRepo{
			DeleteFunc:                 tc.HandleDeleteMatchFunc,
			FindMatchForTournamentFunc: tc.HandleFindMatchForTournamentFunc,
		})
		defer repo.SetMatchRepo(nil)

		repo.SetTournamentRepo(repo.MockTournamentRepo{
			GetFunc: tc.HandleGetTournemantFunc,
		})
		defer repo.SetTournamentRepo(nil)

		req, err := http.NewRequest(http.MethodDelete, "/tournaments/:id/matches/:match_id", nil)
		req = mux.SetURLVars(req, map[string]string{"id": tc.ID, "match_id": tc.MatchID})
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		HandleDeleteMatch(res, req)

		assert.Equal(t, tc.ExpectedStatusCode, res.Code)

		if res.Code == http.StatusNoContent {
			assert.NoError(t, err)
		}
	}
}
