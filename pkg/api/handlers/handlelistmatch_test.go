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
)

func mockListMatchFunc(ctx context.Context) ([]match.Match, errs.AppError) {
	matchMock := prototype.PrototypeMatch()
	matchMock2 := prototype.PrototypeMatch()

	matchMockList := []match.Match{matchMock, matchMock2}

	return matchMockList, nil
}

func mockListMatchThrowFunc(ctx context.Context) ([]match.Match, errs.AppError) {
	return nil, errs.ErrRepoMockAction
}

func TestHandleListMatch(t *testing.T) {
	testCases := []struct {
		Name                string
		HandleListMatchFunc func(ctx context.Context) ([]match.Match, errs.AppError)
		MarshalFunc         func(v interface{}) ([]byte, error)
		WriteFunc           func(http.ResponseWriter, []byte) (int, error)
		ExpectedStatusCode  int
	}{
		{
			Name:                "Success handle list matches",
			HandleListMatchFunc: mockListMatchFunc,
			MarshalFunc:         jsonMarshal,
			WriteFunc:           write,
			ExpectedStatusCode:  200,
		}, {
			Name:                "Throwing handle list matches",
			HandleListMatchFunc: mockListMatchThrowFunc,
			MarshalFunc:         jsonMarshal,
			WriteFunc:           write,
			ExpectedStatusCode:  500,
		}, {
			Name:                "Throwing error on marshal function",
			HandleListMatchFunc: mockListMatchFunc,
			MarshalFunc:         fakeMarshal,
			WriteFunc:           write,
			ExpectedStatusCode:  500,
		}, {
			Name:                "Throwing error on write function",
			HandleListMatchFunc: mockListMatchFunc,
			MarshalFunc:         jsonMarshal,
			WriteFunc:           fakeWrite,
			ExpectedStatusCode:  500,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		repo.SetMatchRepo(repo.MockMatchRepo{
			ListFunc: tc.HandleListMatchFunc,
		})
		defer repo.SetMatchRepo(nil)

		jsonMarshal = tc.MarshalFunc
		defer restoreMarshal(jsonMarshal)

		write = tc.WriteFunc
		defer restoreWrite(write)

		req, err := http.NewRequest(http.MethodGet, "/tournaments/{id}/matches", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "any"})
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		HandleListMatch(res, req)

		assert.Equal(t, tc.ExpectedStatusCode, res.Code)
		t.Logf("Response Body: %v", res.Body)

		if res.Code == http.StatusOK {
			matches := []match.Match{}
			err = json.Unmarshal(res.Body.Bytes(), &matches)
			assert.NoError(t, err)

			assert.Equal(t, 2, len(matches))
		}
	}
}
