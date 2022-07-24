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
	"github.com/rafaelsanzio/go-flashscore/pkg/player"
	"github.com/rafaelsanzio/go-flashscore/pkg/prototype"
)

func mockListPlayerFunc(ctx context.Context) ([]player.Player, errs.AppError) {
	playerMock := prototype.PrototypePlayer()

	playerMock2 := prototype.PrototypePlayer()
	playerMock2.Name = "Real Madrid B"

	playerMockList := []player.Player{playerMock, playerMock2}

	return playerMockList, nil
}

func mockListPlayerThrowFunc(ctx context.Context) ([]player.Player, errs.AppError) {
	return nil, errs.ErrRepoMockAction
}

func TestHandleListPlayer(t *testing.T) {
	testCases := []struct {
		Name                 string
		HandleListPlayerFunc func(ctx context.Context) ([]player.Player, errs.AppError)
		MarshalFunc          func(v interface{}) ([]byte, error)
		WriteFunc            func(http.ResponseWriter, []byte) (int, error)
		ExpectedStatusCode   int
	}{
		{
			Name:                 "Success handle list players",
			HandleListPlayerFunc: mockListPlayerFunc,
			MarshalFunc:          jsonMarshal,
			WriteFunc:            write,
			ExpectedStatusCode:   200,
		}, {
			Name:                 "Throwing handle list players",
			HandleListPlayerFunc: mockListPlayerThrowFunc,
			MarshalFunc:          jsonMarshal,
			WriteFunc:            write,
			ExpectedStatusCode:   500,
		}, {
			Name:                 "Throwing error on marshal function",
			HandleListPlayerFunc: mockListPlayerFunc,
			MarshalFunc:          fakeMarshal,
			WriteFunc:            write,
			ExpectedStatusCode:   500,
		}, {
			Name:                 "Throwing error on write function",
			HandleListPlayerFunc: mockListPlayerFunc,
			MarshalFunc:          jsonMarshal,
			WriteFunc:            fakeWrite,
			ExpectedStatusCode:   500,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		repo.SetPlayerRepo(repo.MockPlayerRepo{
			ListFunc: tc.HandleListPlayerFunc,
		})
		defer repo.SetPlayerRepo(nil)

		jsonMarshal = tc.MarshalFunc
		defer restoreMarshal(jsonMarshal)

		write = tc.WriteFunc
		defer restoreWrite(write)

		req, err := http.NewRequest(http.MethodGet, "/players", nil)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		HandleListPlayer(res, req)

		assert.Equal(t, tc.ExpectedStatusCode, res.Code)
		t.Logf("Response Body: %v", res.Body)

		if res.Code == http.StatusOK {
			player := []player.Player{}
			err = json.Unmarshal(res.Body.Bytes(), &player)
			assert.NoError(t, err)

			assert.Equal(t, 2, len(player))
		}
	}
}
