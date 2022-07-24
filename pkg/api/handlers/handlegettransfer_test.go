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
	"github.com/rafaelsanzio/go-flashscore/pkg/transfer"
)

func mockGetTransferFunc(ctx context.Context, id string) (*transfer.Transfer, errs.AppError) {
	transferMock := prototype.PrototypeTransfer()
	return &transferMock, nil
}

func mockGetTransferThrowFunc(ctx context.Context, id string) (*transfer.Transfer, errs.AppError) {
	return nil, errs.ErrRepoMockAction
}

func mockGetTransferNilFunc(ctx context.Context, id string) (*transfer.Transfer, errs.AppError) {
	return nil, nil
}

func TestHandleGetTransfer(t *testing.T) {
	testCases := []struct {
		Name                  string
		ID                    string
		HandleGetTransferFunc func(ctx context.Context, id string) (*transfer.Transfer, errs.AppError)
		MarshalFunc           func(v interface{}) ([]byte, error)
		WriteFunc             func(http.ResponseWriter, []byte) (int, error)
		ExpectedStatusCode    int
	}{
		{
			Name:                  "Success handle get transfer",
			ID:                    "1",
			HandleGetTransferFunc: mockGetTransferFunc,
			MarshalFunc:           jsonMarshal,
			WriteFunc:             write,
			ExpectedStatusCode:    200,
		}, {
			Name:                  "Not Found handle get transfer",
			ID:                    "",
			HandleGetTransferFunc: mockGetTransferFunc,
			MarshalFunc:           jsonMarshal,
			WriteFunc:             write,
			ExpectedStatusCode:    404,
		}, {
			Name:                  "Getting error on transfer repo",
			ID:                    "1",
			HandleGetTransferFunc: mockGetTransferThrowFunc,
			MarshalFunc:           jsonMarshal,
			WriteFunc:             write,
			ExpectedStatusCode:    500,
		}, {
			Name:                  "Getting error on marshal function",
			ID:                    "1",
			HandleGetTransferFunc: mockGetTransferFunc,
			MarshalFunc:           fakeMarshal,
			WriteFunc:             write,
			ExpectedStatusCode:    500,
		}, {
			Name:                  "Getting error on write function",
			ID:                    "1",
			HandleGetTransferFunc: mockGetTransferFunc,
			MarshalFunc:           jsonMarshal,
			WriteFunc:             fakeWrite,
			ExpectedStatusCode:    500,
		}, {
			Name:                  "Getting error on get function returning nil",
			ID:                    "1",
			HandleGetTransferFunc: mockGetTransferNilFunc,
			MarshalFunc:           jsonMarshal,
			WriteFunc:             write,
			ExpectedStatusCode:    404,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		repo.SetTransferRepo(repo.MockTransferRepo{
			GetFunc: tc.HandleGetTransferFunc,
		})
		defer repo.SetTransferRepo(nil)

		jsonMarshal = tc.MarshalFunc
		defer restoreMarshal(jsonMarshal)

		write = tc.WriteFunc
		defer restoreWrite(write)

		req, err := http.NewRequest(http.MethodGet, "/transfer/:id", nil)
		req = mux.SetURLVars(req, map[string]string{"id": tc.ID})
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		HandleGetTransfer(res, req)

		assert.Equal(t, tc.ExpectedStatusCode, res.Code)
		t.Logf("Response Body: %v", res.Body)

		if res.Code == http.StatusOK {
			team := transfer.Transfer{}
			err = json.Unmarshal(res.Body.Bytes(), &team)
			assert.NoError(t, err)
		}
	}
}
