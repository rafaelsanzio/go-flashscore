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
	"github.com/rafaelsanzio/go-flashscore/pkg/match"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
	"github.com/rafaelsanzio/go-flashscore/pkg/prototype"
	"github.com/rafaelsanzio/go-flashscore/pkg/tournament"
)

func mockFindMatchForTournamentFunc(ctx context.Context, id, tournamentID string) (*match.Match, errs.AppError) {
	matchMock := prototype.PrototypeMatch()
	return &matchMock, nil
}

func mockFindMatchForTournamentThrowFunc(ctx context.Context, id, tournamentID string) (*match.Match, errs.AppError) {
	return nil, errs.ErrRepoMockAction
}

func mockFindMatchForTournamentNilFunc(ctx context.Context, id, tournamentID string) (*match.Match, errs.AppError) {
	return nil, nil
}

func TestHandleGetMatch(t *testing.T) {
	testCases := []struct {
		Name                             string
		ID                               string
		MatchID                          string
		HandleFindMatchForTournamentFunc func(ctx context.Context, id, tournamentID string) (*match.Match, errs.AppError)
		HandleGetTournamentFunc          func(ctx context.Context, id string) (*tournament.Tournament, errs.AppError)
		MarshalFunc                      func(v interface{}) ([]byte, error)
		WriteFunc                        func(http.ResponseWriter, []byte) (int, error)
		ExpectedStatusCode               int
	}{
		{
			Name:                             "Success handle get match",
			ID:                               "1",
			MatchID:                          "1",
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			MarshalFunc:                      jsonMarshal,
			WriteFunc:                        write,
			ExpectedStatusCode:               200,
		}, {
			Name:                             "Not Found id param to handle get match",
			ID:                               "",
			MatchID:                          "1",
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			MarshalFunc:                      jsonMarshal,
			WriteFunc:                        write,
			ExpectedStatusCode:               404,
		}, {
			Name:                             "Not Found match_id param to handle get match",
			ID:                               "1",
			MatchID:                          "",
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			MarshalFunc:                      jsonMarshal,
			WriteFunc:                        write,
			ExpectedStatusCode:               404,
		}, {
			Name:                             "Getting error on get match function",
			ID:                               "1",
			MatchID:                          "1",
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentThrowFunc,
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			MarshalFunc:                      jsonMarshal,
			WriteFunc:                        write,
			ExpectedStatusCode:               500,
		}, {
			Name:                             "Getting error on get func returning nil",
			ID:                               "1",
			MatchID:                          "1",
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentNilFunc,
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			MarshalFunc:                      jsonMarshal,
			WriteFunc:                        write,
			ExpectedStatusCode:               404,
		}, {
			Name:                             "Getting error on get tournament function",
			ID:                               "1",
			MatchID:                          "1",
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentThrowFunc,
			HandleGetTournamentFunc:          mockGetTournamentThrowFunc,
			MarshalFunc:                      jsonMarshal,
			WriteFunc:                        write,
			ExpectedStatusCode:               500,
		}, {
			Name:                             "Getting error on get tournament function retuning nil",
			ID:                               "1",
			MatchID:                          "1",
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentThrowFunc,
			HandleGetTournamentFunc:          mockGetTournamentNilFunc,
			MarshalFunc:                      jsonMarshal,
			WriteFunc:                        write,
			ExpectedStatusCode:               404,
		}, {
			Name:                             "Getting error on marshal function",
			ID:                               "1",
			MatchID:                          "1",
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			MarshalFunc:                      fakeMarshal,
			WriteFunc:                        write,
			ExpectedStatusCode:               500,
		}, {
			Name:                             "Getting error on write function",
			ID:                               "1",
			MatchID:                          "1",
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			MarshalFunc:                      jsonMarshal,
			WriteFunc:                        fakeWrite,
			ExpectedStatusCode:               500,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		repo.SetMatchRepo(repo.MockMatchRepo{
			FindMatchForTournamentFunc: tc.HandleFindMatchForTournamentFunc,
		})
		defer repo.SetMatchRepo(nil)

		repo.SetTournamentRepo(repo.MockTournamentRepo{
			GetFunc: tc.HandleGetTournamentFunc,
		})
		defer repo.SetTournamentRepo(nil)

		jsonMarshal = tc.MarshalFunc
		defer restoreMarshal(jsonMarshal)

		write = tc.WriteFunc
		defer restoreWrite(write)

		req, err := http.NewRequest(http.MethodGet, "/tournaments/:id/matches/:match_id", nil)
		req = mux.SetURLVars(req, map[string]string{"id": tc.ID, "match_id": tc.MatchID})
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		HandleGetMatch(res, req)

		assert.Equal(t, tc.ExpectedStatusCode, res.Code)
		t.Logf("Response Body: %v", res.Body)

		if res.Code == http.StatusOK {
			match := match.Match{}
			err = json.Unmarshal(res.Body.Bytes(), &match)
			assert.NoError(t, err)
		}
	}
}
