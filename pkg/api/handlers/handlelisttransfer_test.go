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
	"github.com/rafaelsanzio/go-flashscore/pkg/transfer"
)

func mockListTransferFunc(ctx context.Context) ([]transfer.Transfer, errs.AppError) {
	tranbsferMock := prototype.PrototypeTransfer()

	tranbsferMock2 := prototype.PrototypeTransfer()

	tranbsferMockList := []transfer.Transfer{tranbsferMock, tranbsferMock2}

	return tranbsferMockList, nil
}

func mockListTransferThrowFunc(ctx context.Context) ([]transfer.Transfer, errs.AppError) {
	return nil, errs.ErrRepoMockAction
}

func TestHandleListTransfer(t *testing.T) {
	testCases := []struct {
		Name                   string
		HandleListTransferFunc func(ctx context.Context) ([]transfer.Transfer, errs.AppError)
		MarshalFunc            func(v interface{}) ([]byte, error)
		WriteFunc              func(http.ResponseWriter, []byte) (int, error)
		ExpectedStatusCode     int
	}{
		{
			Name:                   "Success handle list transfers",
			HandleListTransferFunc: mockListTransferFunc,
			MarshalFunc:            jsonMarshal,
			WriteFunc:              write,
			ExpectedStatusCode:     200,
		}, {
			Name:                   "Throwing handle list transfers",
			HandleListTransferFunc: mockListTransferThrowFunc,
			MarshalFunc:            jsonMarshal,
			WriteFunc:              write,
			ExpectedStatusCode:     500,
		}, {
			Name:                   "Throwing error on marshal function",
			HandleListTransferFunc: mockListTransferFunc,
			MarshalFunc:            fakeMarshal,
			WriteFunc:              write,
			ExpectedStatusCode:     500,
		}, {
			Name:                   "Throwing error on write function",
			HandleListTransferFunc: mockListTransferFunc,
			MarshalFunc:            jsonMarshal,
			WriteFunc:              fakeWrite,
			ExpectedStatusCode:     500,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		repo.SetTransferRepo(repo.MockTransferRepo{
			ListFunc: tc.HandleListTransferFunc,
		})
		defer repo.SetTransferRepo(nil)

		jsonMarshal = tc.MarshalFunc
		defer restoreMarshal(jsonMarshal)

		write = tc.WriteFunc
		defer restoreWrite(write)

		req, err := http.NewRequest(http.MethodGet, "/transfers", nil)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		HandleListTransfer(res, req)

		assert.Equal(t, tc.ExpectedStatusCode, res.Code)
		t.Logf("Response Body: %v", res.Body)

		if res.Code == http.StatusOK {
			player := []transfer.Transfer{}
			err = json.Unmarshal(res.Body.Bytes(), &player)
			assert.NoError(t, err)

			assert.Equal(t, 2, len(player))
		}
	}
}
