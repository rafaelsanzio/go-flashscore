package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
	"github.com/rafaelsanzio/go-flashscore/pkg/money"
	"github.com/rafaelsanzio/go-flashscore/pkg/player"
	"github.com/rafaelsanzio/go-flashscore/pkg/prototype"
	"github.com/rafaelsanzio/go-flashscore/pkg/team"
	"github.com/rafaelsanzio/go-flashscore/pkg/transfer"
)

func mockPostTransferFunc(ctx context.Context, t transfer.Transfer) errs.AppError {
	return nil
}

func mockPostTransferThrowFunc(ctx context.Context, t transfer.Transfer) errs.AppError {
	return errs.ErrRepoMockAction
}

func mockGetTeamFuncForTransfer(ctx context.Context, id string) (*team.Team, errs.AppError) {
	teamMock := prototype.PrototypeTeam()
	if id == "any_team_destiny_id" {
		teamMock.ID = "2"
	}
	return &teamMock, nil
}

func TestHandlePostTransfer(t *testing.T) {
	body, err := json.Marshal(TransferEntityPayload{
		Player:         "any_player_id",
		TeamOrigin:     "any_team_origin_id",
		TeamDestiny:    "any_team_destiny_id",
		Amount:         money.Money{Cents: 1000, Currency: money.USD},
		DateOfTransfer: "1990-01-01",
	})
	assert.Equal(t, nil, err)

	goodReq := httptest.NewRequest(http.MethodPost, "/transfers", nil)
	goodReq = mux.SetURLVars(goodReq, map[string]string{})
	goodReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	noBodyReq := httptest.NewRequest(http.MethodPost, "/transfers", nil)
	noBodyReq = mux.SetURLVars(noBodyReq, map[string]string{})

	throwReq := httptest.NewRequest(http.MethodPost, "/transfers", nil)
	throwReq = mux.SetURLVars(throwReq, map[string]string{})
	throwReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	goodReq2 := httptest.NewRequest(http.MethodPost, "/transfers", nil)
	goodReq2 = mux.SetURLVars(goodReq2, map[string]string{})
	goodReq2.Body = ioutil.NopCloser(bytes.NewReader(body))

	testCases := []struct {
		Name                  string
		Request               *http.Request
		HandlePostFunc        func(ctx context.Context, p transfer.Transfer) errs.AppError
		HandleGetTeamFunc     func(ctx context.Context, id string) (*team.Team, errs.AppError)
		HandleGetPlayerFunc   func(ctx context.Context, id string) (*player.Player, errs.AppError)
		ConvertingPayloadFunc func(ctx context.Context, p TransferEntityPayload) (*transfer.Transfer, errs.AppError)
		ExpectedStatusCode    int
	}{
		{
			Name:                  "Should return 201 if successful",
			Request:               goodReq,
			HandlePostFunc:        mockPostTransferFunc,
			HandleGetTeamFunc:     mockGetTeamFuncForTransfer,
			HandleGetPlayerFunc:   mockGetPlayerFunc,
			ConvertingPayloadFunc: convertAndValidatePayloadToTransferFunc,
			ExpectedStatusCode:    201,
		}, {
			Name:                  "Should return 422 bad request",
			Request:               noBodyReq,
			HandlePostFunc:        mockPostTransferFunc,
			HandleGetTeamFunc:     mockGetTeamFuncForTransfer,
			HandleGetPlayerFunc:   mockGetPlayerFunc,
			ConvertingPayloadFunc: convertAndValidatePayloadToTransferFunc,
			ExpectedStatusCode:    422,
		}, {
			Name:                  "Should return 422 throwing error on converting payload func",
			Request:               goodReq2,
			HandlePostFunc:        mockPostTransferFunc,
			HandleGetTeamFunc:     mockGetTeamFuncForTransfer,
			HandleGetPlayerFunc:   mockGetPlayerFunc,
			ConvertingPayloadFunc: fakeConvertAndValidatePayloadToTransfer,
			ExpectedStatusCode:    422,
		}, {
			Name:                  "Should return 500 throwing error on function",
			Request:               throwReq,
			HandlePostFunc:        mockPostTransferThrowFunc,
			HandleGetTeamFunc:     mockGetTeamFuncForTransfer,
			HandleGetPlayerFunc:   mockGetPlayerFunc,
			ConvertingPayloadFunc: convertAndValidatePayloadToTransferFunc,
			ExpectedStatusCode:    500,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		repo.SetTransferRepo(repo.MockTransferRepo{
			InsertFunc: tc.HandlePostFunc,
		})
		defer repo.SetTransferRepo(nil)

		repo.SetTeamRepo(repo.MockTeamRepo{
			GetFunc: tc.HandleGetTeamFunc,
		})
		defer repo.SetTeamRepo(nil)

		repo.SetPlayerRepo(repo.MockPlayerRepo{
			GetFunc: tc.HandleGetPlayerFunc,
		})
		defer repo.SetPlayerRepo(nil)

		convertAndValidatePayloadToTransferFunc = tc.ConvertingPayloadFunc
		defer restoreConvertAndValidatePayloadToTransfer(convertAndValidatePayloadToTransferFunc)

		w := httptest.NewRecorder()

		HandlePostTransfer(w, tc.Request)
		assert.Equal(t, tc.ExpectedStatusCode, w.Code)
	}
}

func TestConvertPayloadToTransfer(t *testing.T) {
	inPayload := TransferEntityPayload{
		Player:         "any_player_id",
		TeamOrigin:     "any_team_origin_id",
		TeamDestiny:    "any_team_destiny_id",
		Amount:         money.Money{Cents: 1000, Currency: money.USD},
		DateOfTransfer: "1990-01-01",
	}

	expectedTeamDestiny := prototype.PrototypeTeam()
	expectedTeamDestiny.ID = "2"

	expectedTeam := transfer.Transfer{
		ID:             "",
		Player:         prototype.PrototypePlayer(),
		TeamDestiny:    expectedTeamDestiny,
		Amount:         money.Money{Cents: 1000, Currency: money.USD},
		DateOfTransfer: "1990-01-01",
	}

	testCases := []struct {
		Name                string
		Payload             TransferEntityPayload
		HandleGetTeamFunc   func(ctx context.Context, id string) (*team.Team, errs.AppError)
		HandleGetPlayerFunc func(ctx context.Context, id string) (*player.Player, errs.AppError)
		ParsingTimeFunc     func(layout string, value string) (time.Time, error)
		ExpectedTeam        transfer.Transfer
		ExpectError         bool
		ExpectedError       string
	}{
		{
			Name:                "Test Case: 1 - correct body, no error",
			Payload:             inPayload,
			HandleGetTeamFunc:   mockGetTeamFuncForTransfer,
			HandleGetPlayerFunc: mockGetPlayerFunc,
			ParsingTimeFunc:     timeParse,
			ExpectedTeam:        expectedTeam,
			ExpectError:         false,
		}, {
			Name:                "Test Case: 2 - throwing error on get function",
			Payload:             inPayload,
			HandleGetTeamFunc:   mockGetTeamThrowFunc,
			HandleGetPlayerFunc: mockGetPlayerFunc,
			ParsingTimeFunc:     timeParse,
			ExpectedTeam:        expectedTeam,
			ExpectError:         true,
		}, {
			Name:                "Test Case: 3 - throwing error on get team nil function",
			Payload:             inPayload,
			HandleGetTeamFunc:   mockGetTeamNilFunc,
			HandleGetPlayerFunc: mockGetPlayerFunc,
			ParsingTimeFunc:     timeParse,
			ExpectedTeam:        expectedTeam,
			ExpectError:         true,
		}, {
			Name:                "Test Case: 4 - throwing error on get player nil function",
			Payload:             inPayload,
			HandleGetTeamFunc:   mockGetTeamFuncForTransfer,
			HandleGetPlayerFunc: mockGetPlayerNilFunc,
			ParsingTimeFunc:     timeParse,
			ExpectedTeam:        expectedTeam,
			ExpectError:         true,
		}, {
			Name:                "Test Case: 5 - throwing error on get player function",
			Payload:             inPayload,
			HandleGetTeamFunc:   mockGetTeamFuncForTransfer,
			HandleGetPlayerFunc: mockGetPlayerThrowFunc,
			ParsingTimeFunc:     timeParse,
			ExpectedTeam:        expectedTeam,
			ExpectError:         true,
		}, {
			Name:                "Test Case: 6 - throwing error parsing time function",
			Payload:             inPayload,
			HandleGetTeamFunc:   mockGetTeamFuncForTransfer,
			HandleGetPlayerFunc: mockGetPlayerFunc,
			ParsingTimeFunc:     fakeTimeParse,
			ExpectedTeam:        expectedTeam,
			ExpectError:         true,
		}, {
			Name:                "Test Case: 7 - throwing error validation with same teams",
			Payload:             inPayload,
			HandleGetTeamFunc:   mockGetTeamFunc,
			HandleGetPlayerFunc: mockGetPlayerFunc,
			ParsingTimeFunc:     timeParse,
			ExpectedTeam:        expectedTeam,
			ExpectError:         true,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		repo.SetTeamRepo(repo.MockTeamRepo{
			GetFunc: tc.HandleGetTeamFunc,
		})
		defer repo.SetTeamRepo(nil)

		repo.SetPlayerRepo(repo.MockPlayerRepo{
			GetFunc: tc.HandleGetPlayerFunc,
		})
		defer repo.SetPlayerRepo(nil)

		timeParse = tc.ParsingTimeFunc
		defer restoreTimeParse(timeParse)

		player, err := convertAndValidatePayloadToTransferFunc(context.Background(), tc.Payload)
		if tc.ExpectError {
			assert.NotNil(t, err)
		} else {
			assert.Equal(t, tc.ExpectedTeam, *player)
		}
	}
}

func TestDecodeTransferRequest(t *testing.T) {
	body, err := json.Marshal(TransferEntityPayload{
		Player:         "any_player_id",
		TeamOrigin:     "any_team_origin_id",
		TeamDestiny:    "any_team_destiny_id",
		Amount:         money.Money{Cents: 1000, Currency: money.USD},
		DateOfTransfer: "1990-01-01",
	})
	assert.Equal(t, nil, err)

	goodReq := httptest.NewRequest(http.MethodPost, "/transfers", nil)
	goodReq = mux.SetURLVars(goodReq, map[string]string{})

	goodReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	noBodyReq := httptest.NewRequest(http.MethodPost, "/transfers", nil)
	noBodyReq = mux.SetURLVars(noBodyReq, map[string]string{})

	testCases := []struct {
		Name          string
		Request       *http.Request
		Payload       *TransferEntityPayload
		ExpectedError bool
	}{
		{
			Name:    "Test Case: 1 - correct body, no error",
			Request: goodReq, Payload: &TransferEntityPayload{
				Player:         "any_player_id",
				TeamOrigin:     "any_team_origin_id",
				TeamDestiny:    "any_team_destiny_id",
				Amount:         money.Money{Cents: 1000, Currency: money.USD},
				DateOfTransfer: "1990-01-01",
			}, ExpectedError: false,
		},
		{Name: "Test Case: 2 - no body, error found", Request: noBodyReq, Payload: nil, ExpectedError: true},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		decodedPayload, err := decodeTransferRequest(tc.Request)
		if tc.ExpectedError {
			assert.NotNil(t, err)
		} else {
			assert.Equal(t, *tc.Payload, decodedPayload)
		}
	}
}
