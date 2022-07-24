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
	"github.com/rafaelsanzio/go-flashscore/pkg/match"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
	"github.com/rafaelsanzio/go-flashscore/pkg/prototype"
	"github.com/rafaelsanzio/go-flashscore/pkg/team"
	"github.com/rafaelsanzio/go-flashscore/pkg/tournament"
)

func mockPostMatchFunc(ctx context.Context, mt match.Match) errs.AppError {
	return nil
}

func mockPostMatchThrowFunc(ctx context.Context, mt match.Match) errs.AppError {
	return errs.ErrRepoMockAction
}

func mockGetTeamFuncForMatch(ctx context.Context, id string) (*team.Team, errs.AppError) {
	teamMock := prototype.PrototypeTeam()
	if id == "any_away_team_id" {
		teamMock.ID = "2"
	}
	return &teamMock, nil
}

func mockGetTeamThrownFuncForMatch(ctx context.Context, id string) (*team.Team, errs.AppError) {
	teamMock := prototype.PrototypeTeam()
	if id == "any_away_team_id" {
		return nil, errs.ErrRepoMockAction
	}
	return &teamMock, nil
}

func fakeTimeParseForTimeMatch(layout string, value string) (time.Time, error) {
	if layout == "15:04" {
		return time.Time{}, errs.ErrParsingTime
	}
	return time.Time{}, nil
}

func TestHandlePostMatch(t *testing.T) {
	body, err := json.Marshal(MatchEntityPayload{
		HomeTeam:    "any_home_team_id",
		AwayTeam:    "any_away_team_id",
		DateOfMatch: "2022-01-01",
		TimeOfMatch: "16:00",
	})
	assert.Equal(t, nil, err)

	goodReq := httptest.NewRequest(http.MethodPost, "/tournaments/:id/matches", nil)
	goodReq = mux.SetURLVars(goodReq, map[string]string{"id": "any"})
	goodReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	noBodyReq := httptest.NewRequest(http.MethodPost, "/tournaments/:id/matches", nil)
	noBodyReq = mux.SetURLVars(noBodyReq, map[string]string{"id": "any"})

	throwReq := httptest.NewRequest(http.MethodPost, "/tournaments/:id/matches", nil)
	throwReq = mux.SetURLVars(throwReq, map[string]string{"id": "any"})

	throwReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	goodReq2 := httptest.NewRequest(http.MethodPost, "/tournaments/:id/matches", nil)
	goodReq2 = mux.SetURLVars(goodReq2, map[string]string{"id": "any"})
	goodReq2.Body = ioutil.NopCloser(bytes.NewReader(body))

	emptyIDReq := httptest.NewRequest(http.MethodPost, "/tournaments/:id/matches", nil)
	emptyIDReq = mux.SetURLVars(emptyIDReq, map[string]string{"id": ""})
	emptyIDReq.Body = ioutil.NopCloser(bytes.NewReader(body))

	testCases := []struct {
		Name                    string
		Request                 *http.Request
		HandlePostMatchFunc     func(ctx context.Context, mt match.Match) errs.AppError
		HandleGetTournamentFunc func(ctx context.Context, id string) (*tournament.Tournament, errs.AppError)
		HandleGetTeamFunc       func(ctx context.Context, id string) (*team.Team, errs.AppError)
		ExpectedStatusCode      int
	}{
		{
			Name:                    "Should return 201 if successful",
			Request:                 goodReq,
			HandlePostMatchFunc:     mockPostMatchFunc,
			HandleGetTournamentFunc: mockGetTournamentFunc,
			HandleGetTeamFunc:       mockGetTeamFuncForMatch,
			ExpectedStatusCode:      201,
		}, {
			Name:                    "Should return 422 bad request",
			Request:                 noBodyReq,
			HandlePostMatchFunc:     mockPostMatchThrowFunc,
			HandleGetTournamentFunc: mockGetTournamentFunc,
			HandleGetTeamFunc:       mockGetTeamFunc,
			ExpectedStatusCode:      422,
		}, {
			Name:                    "Should return 500 throwing erro on insert match function",
			Request:                 goodReq2,
			HandlePostMatchFunc:     mockPostMatchThrowFunc,
			HandleGetTournamentFunc: mockGetTournamentFunc,
			HandleGetTeamFunc:       mockGetTeamFuncForMatch,
			ExpectedStatusCode:      500,
		}, {
			Name:                    "Should return 404 with empty id",
			Request:                 emptyIDReq,
			HandlePostMatchFunc:     mockPostMatchFunc,
			HandleGetTournamentFunc: mockGetTournamentFunc,
			HandleGetTeamFunc:       mockGetTeamFunc,
			ExpectedStatusCode:      404,
		}, {
			Name:                    "Should return 422 throwing error on request",
			Request:                 throwReq,
			HandlePostMatchFunc:     mockPostMatchFunc,
			HandleGetTournamentFunc: mockGetTournamentFunc,
			HandleGetTeamFunc:       mockGetTeamFunc,
			ExpectedStatusCode:      422,
		}, {
			Name:                    "Should return 500 throwing error on get tournament function",
			Request:                 goodReq2,
			HandlePostMatchFunc:     mockPostMatchFunc,
			HandleGetTournamentFunc: mockGetTournamentThrowFunc,
			HandleGetTeamFunc:       mockGetTeamFunc,
			ExpectedStatusCode:      500,
		}, {
			Name:                    "Should return 422 throwing error on get team function",
			Request:                 goodReq2,
			HandlePostMatchFunc:     mockPostMatchFunc,
			HandleGetTournamentFunc: mockGetTournamentFunc,
			HandleGetTeamFunc:       mockGetTeamThrowFunc,
			ExpectedStatusCode:      422,
		}, {
			Name:                    "Should return 404 throwing error on get tournament function returning nil",
			Request:                 goodReq2,
			HandlePostMatchFunc:     mockPostMatchFunc,
			HandleGetTournamentFunc: mockGetTournamentNilFunc,
			HandleGetTeamFunc:       mockGetTeamThrowFunc,
			ExpectedStatusCode:      404,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		repo.SetMatchRepo(repo.MockMatchRepo{
			InsertFunc: tc.HandlePostMatchFunc,
		})
		defer repo.SetMatchRepo(nil)

		repo.SetTournamentRepo(repo.MockTournamentRepo{
			GetFunc: tc.HandleGetTournamentFunc,
		})
		defer repo.SetTournamentRepo(nil)

		repo.SetTeamRepo(repo.MockTeamRepo{
			GetFunc: tc.HandleGetTeamFunc,
		})
		defer repo.SetTeamRepo(nil)

		w := httptest.NewRecorder()

		HandlePostMatch(w, tc.Request)
		assert.Equal(t, tc.ExpectedStatusCode, w.Code)
	}
}

func TestConvertAndValidatePayloadToMatch(t *testing.T) {
	inPayload := MatchEntityPayload{
		HomeTeam:    "any_home_team_id",
		AwayTeam:    "any_away_team_id",
		DateOfMatch: "2022-01-01",
		TimeOfMatch: "16:00",
	}

	expectedAwayTeam := prototype.PrototypeTeam()
	expectedAwayTeam.ID = "2"

	expectedMatch := match.Match{
		ID:          "",
		HomeTeam:    prototype.PrototypeTeam(),
		AwayTeam:    expectedAwayTeam,
		DateOfMatch: "2022-01-01",
		TimeOfMatch: "16:00",
	}

	testCases := []struct {
		Name              string
		Payload           MatchEntityPayload
		HandleGetTeamFunc func(ctx context.Context, id string) (*team.Team, errs.AppError)
		ParsingTimeFunc   func(layout string, value string) (time.Time, error)
		ExpectedMatch     match.Match
		ExpectError       bool
		ExpectedError     string
	}{
		{
			Name:              "Test Case: 1 - correct body, no error",
			Payload:           inPayload,
			HandleGetTeamFunc: mockGetTeamFuncForMatch,
			ParsingTimeFunc:   timeParse,
			ExpectedMatch:     expectedMatch,
			ExpectError:       false,
		}, {
			Name:              "Test Case: 2 - throwing error on get function",
			Payload:           inPayload,
			HandleGetTeamFunc: mockGetTeamThrowFunc,
			ParsingTimeFunc:   timeParse,
			ExpectedMatch:     expectedMatch,
			ExpectError:       true,
		}, {
			Name:              "Test Case: 3 - throwing error on get function",
			Payload:           inPayload,
			HandleGetTeamFunc: mockGetTeamThrownFuncForMatch,
			ParsingTimeFunc:   timeParse,
			ExpectedMatch:     expectedMatch,
			ExpectError:       true,
		}, {
			Name:              "Test Case: 4 - throwing error on get team nil function",
			Payload:           inPayload,
			HandleGetTeamFunc: mockGetTeamNilFunc,
			ParsingTimeFunc:   timeParse,
			ExpectedMatch:     expectedMatch,
			ExpectError:       true,
		}, {
			Name:              "Test Case: 5 - throwing error parsing time function",
			Payload:           inPayload,
			HandleGetTeamFunc: mockGetTeamFuncForMatch,
			ParsingTimeFunc:   fakeTimeParse,
			ExpectedMatch:     expectedMatch,
			ExpectError:       true,
		}, {
			Name:              "Test Case: 6 - throwing error parsing time function",
			Payload:           inPayload,
			HandleGetTeamFunc: mockGetTeamFuncForMatch,
			ParsingTimeFunc:   fakeTimeParseForTimeMatch,
			ExpectedMatch:     expectedMatch,
			ExpectError:       true,
		}, {
			Name:              "Test Case: 7 - throwing error validation with same teams",
			Payload:           inPayload,
			HandleGetTeamFunc: mockGetTeamFunc,
			ParsingTimeFunc:   timeParse,
			ExpectedMatch:     expectedMatch,
			ExpectError:       true,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.Name)

		repo.SetTeamRepo(repo.MockTeamRepo{
			GetFunc: tc.HandleGetTeamFunc,
		})
		defer repo.SetTeamRepo(nil)

		timeParse = tc.ParsingTimeFunc
		defer restoreTimeParse(timeParse)

		match, err := convertAndValidatePayloadToMatch(context.Background(), tc.Payload)
		if tc.ExpectError {
			assert.NotNil(t, err)
		} else {
			assert.Equal(t, tc.ExpectedMatch, *match)
		}
	}
}
