package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/event"
	"github.com/rafaelsanzio/go-flashscore/pkg/match"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
	"github.com/rafaelsanzio/go-flashscore/pkg/tournament"
)

func TestHandlePostMatchHalftime(t *testing.T) {

	goodReq := httptest.NewRequest(http.MethodPost, "/tournament/{id}/matches/{match_id}/events/halftime", nil)
	goodReq = mux.SetURLVars(goodReq, map[string]string{"id": "any", "match_id": "any"})

	missParamIDReq := httptest.NewRequest(http.MethodPost, "/tournament/{id}/matches/{match_id}/events/halftime", nil)
	missParamIDReq = mux.SetURLVars(missParamIDReq, map[string]string{"id": "", "match_id": "any"})

	throwReq := httptest.NewRequest(http.MethodPost, "/tournament/{id}/matches/{match_id}/events/halftime", nil)
	throwReq = mux.SetURLVars(throwReq, map[string]string{"id": "any", "match_id": "any"})

	goodReq2 := httptest.NewRequest(http.MethodPost, "/tournament/{id}/matches/{match_id}/events/halftime", nil)
	mux.SetURLVars(goodReq2, map[string]string{"id": "any", "match_id": "any"})

	missParamMatchIDReq := httptest.NewRequest(http.MethodPost, "/tournament/{id}/matches/{match_id}/events/halftime", nil)
	missParamMatchIDReq = mux.SetURLVars(missParamMatchIDReq, map[string]string{"id": "any", "match_id": ""})

	testCases := []struct {
		Name                             string
		Request                          *http.Request
		HandleGetTournamentFunc          func(ctx context.Context, id string) (*tournament.Tournament, errs.AppError)
		HandleFindMatchForTournamentFunc func(ctx context.Context, id, tournamentID string) (*match.Match, errs.AppError)
		HandlePostEventFunc              func(ctx context.Context, e event.Event) errs.AppError
		ExpectedStatusCode               int
	}{
		{
			Name:                             "Should return 201 if successful",
			Request:                          goodReq,
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchStatusInProgressForTournamentFunc,
			HandlePostEventFunc:              mockPostEventFunc,
			ExpectedStatusCode:               201,
		}, {
			Name:                             "Should return 404 missing id param",
			Request:                          missParamIDReq,
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchStatusInProgressForTournamentFunc,
			HandlePostEventFunc:              mockPostEventFunc,
			ExpectedStatusCode:               404,
		}, {
			Name:                             "Should return 500 throwing error get tournament function",
			Request:                          throwReq,
			HandleGetTournamentFunc:          mockGetTournamentThrowFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchStatusInProgressForTournamentFunc,
			HandlePostEventFunc:              mockPostEventFunc,
			ExpectedStatusCode:               500,
		}, {
			Name:                             "Should return 404 if tournament is not found",
			Request:                          throwReq,
			HandleGetTournamentFunc:          mockGetTournamentNilFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchStatusInProgressForTournamentFunc,
			HandlePostEventFunc:              mockPostEventFunc,
			ExpectedStatusCode:               404,
		}, {
			Name:                             "Should return 404 missing match id param",
			Request:                          missParamMatchIDReq,
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchStatusInProgressForTournamentFunc,
			HandlePostEventFunc:              mockPostEventFunc,
			ExpectedStatusCode:               404,
		}, {
			Name:                             "Should return 500 throwing error find match to tournament function",
			Request:                          throwReq,
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentThrowFunc,
			HandlePostEventFunc:              mockPostEventFunc,
			ExpectedStatusCode:               500,
		}, {
			Name:                             "Should return 404 if match is not found",
			Request:                          throwReq,
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentNilFunc,
			HandlePostEventFunc:              mockPostEventFunc,
			ExpectedStatusCode:               404,
		}, {
			Name:                             "Should return 422 if match is not started",
			Request:                          throwReq,
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchStatusNotStartedForTournamentFunc,
			HandlePostEventFunc:              mockPostEventFunc,
			ExpectedStatusCode:               422,
		}, {
			Name:                             "Should return 500 throwing error post event function",
			Request:                          throwReq,
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchStatusInProgressForTournamentFunc,
			HandlePostEventFunc:              mockPostEventThrowFunc,
			ExpectedStatusCode:               500,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		repo.SetTournamentRepo(repo.MockTournamentRepo{
			GetFunc: tc.HandleGetTournamentFunc,
		})
		defer repo.SetTournamentRepo(nil)

		repo.SetMatchRepo(repo.MockMatchRepo{
			FindMatchForTournamentFunc: tc.HandleFindMatchForTournamentFunc,
		})
		defer repo.SetMatchRepo(nil)

		repo.SetEventRepo(repo.MockEventRepo{
			InsertFunc: tc.HandlePostEventFunc,
		})
		defer repo.SetEventRepo(nil)

		w := httptest.NewRecorder()

		HandlePostMatchHalftime(w, tc.Request)
		assert.Equal(t, tc.ExpectedStatusCode, w.Code)
	}
}
