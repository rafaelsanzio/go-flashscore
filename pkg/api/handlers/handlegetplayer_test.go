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
	"github.com/rafaelsanzio/go-flashscore/pkg/player"
	"github.com/rafaelsanzio/go-flashscore/pkg/prototype"
)

func mockGetPlayerFunc(ctx context.Context, id string) (*player.Player, errs.AppError) {
	playerMock := prototype.PrototypePlayer()
	return &playerMock, nil
}

func mockGetPlayerThrowFunc(ctx context.Context, id string) (*player.Player, errs.AppError) {
	return nil, errs.ErrRepoMockAction
}

func mockGetPlayerNilFunc(ctx context.Context, id string) (*player.Player, errs.AppError) {
	return nil, nil
}

func TestHandleGetPlayer(t *testing.T) {
	testCases := []struct {
		Name                string
		ID                  string
		HandleGetPlayerFunc func(ctx context.Context, id string) (*player.Player, errs.AppError)
		MarshalFunc         func(v interface{}) ([]byte, error)
		WriteFunc           func(http.ResponseWriter, []byte) (int, error)
		ExpectedStatusCode  int
	}{
		{
			Name:                "Success handle get player",
			ID:                  "1",
			HandleGetPlayerFunc: mockGetPlayerFunc,
			MarshalFunc:         jsonMarshal,
			WriteFunc:           write,
			ExpectedStatusCode:  200,
		}, {
			Name:                "Not Found handle get player",
			ID:                  "",
			HandleGetPlayerFunc: mockGetPlayerFunc,
			MarshalFunc:         jsonMarshal,
			WriteFunc:           write,
			ExpectedStatusCode:  404,
		}, {
			Name:                "Getting error on player repo",
			ID:                  "1",
			HandleGetPlayerFunc: mockGetPlayerThrowFunc,
			MarshalFunc:         jsonMarshal,
			WriteFunc:           write,
			ExpectedStatusCode:  500,
		}, {
			Name:                "Getting error on marshal function",
			ID:                  "1",
			HandleGetPlayerFunc: mockGetPlayerFunc,
			MarshalFunc:         fakeMarshal,
			WriteFunc:           write,
			ExpectedStatusCode:  500,
		}, {
			Name:                "Getting error on write function",
			ID:                  "1",
			HandleGetPlayerFunc: mockGetPlayerFunc,
			MarshalFunc:         jsonMarshal,
			WriteFunc:           fakeWrite,
			ExpectedStatusCode:  500,
		}, {
			Name:                "Getting error on get func returning nil",
			ID:                  "1",
			HandleGetPlayerFunc: mockGetPlayerNilFunc,
			MarshalFunc:         jsonMarshal,
			WriteFunc:           write,
			ExpectedStatusCode:  404,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		repo.SetPlayerRepo(repo.MockPlayerRepo{
			GetFunc: tc.HandleGetPlayerFunc,
		})
		defer repo.SetPlayerRepo(nil)

		jsonMarshal = tc.MarshalFunc
		defer restoreMarshal(jsonMarshal)

		write = tc.WriteFunc
		defer restoreWrite(write)

		req, err := http.NewRequest(http.MethodGet, "/players/:id", nil)
		req = mux.SetURLVars(req, map[string]string{"id": tc.ID})
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		HandleGetPlayer(res, req)

		assert.Equal(t, tc.ExpectedStatusCode, res.Code)
		t.Logf("Response Body: %v", res.Body)

		if res.Code == http.StatusOK {
			team := player.Player{}
			err = json.Unmarshal(res.Body.Bytes(), &team)
			assert.NoError(t, err)
		}
	}
}
