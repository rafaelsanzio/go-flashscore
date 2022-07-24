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
	"github.com/rafaelsanzio/go-flashscore/pkg/tournament"
)

func mockListTournamentFunc(ctx context.Context) ([]tournament.Tournament, errs.AppError) {
	tournamentMock := prototype.PrototypeTournament()
	tournamentMock2 := prototype.PrototypeTournament()

	tournamentMockList := []tournament.Tournament{tournamentMock, tournamentMock2}

	return tournamentMockList, nil
}

func mockListTournamentThrowFunc(ctx context.Context) ([]tournament.Tournament, errs.AppError) {
	return nil, errs.ErrRepoMockAction
}

func TestHandleListTournament(t *testing.T) {
	testCases := []struct {
		Name                     string
		HandleListTournamentFunc func(ctx context.Context) ([]tournament.Tournament, errs.AppError)
		MarshalFunc              func(v interface{}) ([]byte, error)
		WriteFunc                func(http.ResponseWriter, []byte) (int, error)
		ExpectedStatusCode       int
	}{
		{
			Name:                     "Success handle list tournaments",
			HandleListTournamentFunc: mockListTournamentFunc,
			MarshalFunc:              jsonMarshal,
			WriteFunc:                write,
			ExpectedStatusCode:       200,
		}, {
			Name:                     "Throwing handle list tournaments",
			HandleListTournamentFunc: mockListTournamentThrowFunc,
			MarshalFunc:              jsonMarshal,
			WriteFunc:                write,
			ExpectedStatusCode:       500,
		}, {
			Name:                     "Throwing error on marshal function",
			HandleListTournamentFunc: mockListTournamentFunc,
			MarshalFunc:              fakeMarshal,
			WriteFunc:                write,
			ExpectedStatusCode:       500,
		}, {
			Name:                     "Throwing error on write function",
			HandleListTournamentFunc: mockListTournamentFunc,
			MarshalFunc:              jsonMarshal,
			WriteFunc:                fakeWrite,
			ExpectedStatusCode:       500,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		repo.SetTournamentRepo(repo.MockTournamentRepo{
			ListFunc: tc.HandleListTournamentFunc,
		})
		defer repo.SetTournamentRepo(nil)

		jsonMarshal = tc.MarshalFunc
		defer restoreMarshal(jsonMarshal)

		write = tc.WriteFunc
		defer restoreWrite(write)

		req, err := http.NewRequest(http.MethodGet, "/tournaments", nil)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		HandleListTournament(res, req)

		assert.Equal(t, tc.ExpectedStatusCode, res.Code)
		t.Logf("Response Body: %v", res.Body)

		if res.Code == http.StatusOK {
			tournament := []tournament.Tournament{}
			err = json.Unmarshal(res.Body.Bytes(), &tournament)
			assert.NoError(t, err)

			assert.Equal(t, 2, len(tournament))
		}
	}
}
