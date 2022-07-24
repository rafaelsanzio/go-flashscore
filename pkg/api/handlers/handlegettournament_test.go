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
	"github.com/rafaelsanzio/go-flashscore/pkg/tournament"
)

func mockGetTournamentFunc(ctx context.Context, id string) (*tournament.Tournament, errs.AppError) {
	tournamentMock := prototype.PrototypeTournament()
	return &tournamentMock, nil
}

func mockGetTournamentThrowFunc(ctx context.Context, id string) (*tournament.Tournament, errs.AppError) {
	return nil, errs.ErrRepoMockAction
}

func mockGetTournamentNilFunc(ctx context.Context, id string) (*tournament.Tournament, errs.AppError) {
	return nil, nil
}

func TestHandleGetTournament(t *testing.T) {
	testCases := []struct {
		Name                    string
		ID                      string
		HandleGetTournamentFunc func(ctx context.Context, id string) (*tournament.Tournament, errs.AppError)
		MarshalFunc             func(v interface{}) ([]byte, error)
		WriteFunc               func(http.ResponseWriter, []byte) (int, error)
		ExpectedStatusCode      int
	}{
		{
			Name:                    "Success handle get tournament",
			ID:                      "1",
			HandleGetTournamentFunc: mockGetTournamentFunc,
			MarshalFunc:             jsonMarshal,
			WriteFunc:               write,
			ExpectedStatusCode:      200,
		}, {
			Name:                    "Not Found handle get tournament",
			ID:                      "",
			HandleGetTournamentFunc: mockGetTournamentFunc,
			MarshalFunc:             jsonMarshal,
			WriteFunc:               write,
			ExpectedStatusCode:      404,
		}, {
			Name:                    "Getting error on tournament repo",
			ID:                      "1",
			HandleGetTournamentFunc: mockGetTournamentThrowFunc,
			MarshalFunc:             jsonMarshal,
			WriteFunc:               write,
			ExpectedStatusCode:      500,
		}, {
			Name:                    "Getting error on marshal function",
			ID:                      "1",
			HandleGetTournamentFunc: mockGetTournamentFunc,
			MarshalFunc:             fakeMarshal,
			WriteFunc:               write,
			ExpectedStatusCode:      500,
		}, {
			Name:                    "Getting error on write function",
			ID:                      "1",
			HandleGetTournamentFunc: mockGetTournamentFunc,
			MarshalFunc:             jsonMarshal,
			WriteFunc:               fakeWrite,
			ExpectedStatusCode:      500,
		}, {
			Name:                    "Getting error on get func returning nil",
			ID:                      "1",
			HandleGetTournamentFunc: mockGetTournamentNilFunc,
			MarshalFunc:             jsonMarshal,
			WriteFunc:               write,
			ExpectedStatusCode:      404,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		repo.SetTournamentRepo(repo.MockTournamentRepo{
			GetFunc: tc.HandleGetTournamentFunc,
		})
		defer repo.SetTournamentRepo(nil)

		jsonMarshal = tc.MarshalFunc
		defer restoreMarshal(jsonMarshal)

		write = tc.WriteFunc
		defer restoreWrite(write)

		req, err := http.NewRequest(http.MethodGet, "/tournaments/:id", nil)
		req = mux.SetURLVars(req, map[string]string{"id": tc.ID})
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		HandleGetTournament(res, req)

		assert.Equal(t, tc.ExpectedStatusCode, res.Code)
		t.Logf("Response Body: %v", res.Body)

		if res.Code == http.StatusOK {
			tournament := tournament.Tournament{}
			err = json.Unmarshal(res.Body.Bytes(), &tournament)
			assert.NoError(t, err)
		}
	}
}
