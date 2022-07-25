package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/player"
	"github.com/rafaelsanzio/go-flashscore/pkg/team"
	"github.com/rafaelsanzio/go-flashscore/pkg/tournament"
	"github.com/rafaelsanzio/go-flashscore/pkg/transfer"
)

var rateLimitAllow = rateLimit.Allow

func fakeRateLimitAllow() bool {
	return false
}

func restoreRateLimitAllow(replace func() bool) {
	rateLimitAllow = replace
}

var jsonMarshal = json.Marshal

func fakeMarshal(v interface{}) ([]byte, error) {
	return []byte{}, errs.ErrMarshalingJson
}

func restoreMarshal(replace func(v interface{}) ([]byte, error)) {
	jsonMarshal = replace
}

var write = http.ResponseWriter.Write

func fakeWrite(http.ResponseWriter, []byte) (int, error) {
	return 0, errs.ErrResponseWriter
}

func restoreWrite(replace func(http.ResponseWriter, []byte) (int, error)) {
	write = replace
}

var convertPayloadToTeamFunc = convertPayloadToTeam

func fakeConvertPayloadToTeamFunc(t TeamEntityPayload) (team.Team, errs.AppError) {
	return team.Team{}, errs.ErrConvertingPayload
}

func restoreConvertPayloadToTeamFunc(replace func(t TeamEntityPayload) (team.Team, errs.AppError)) {
	convertPayloadToTeamFunc = replace
}

var convertPayloadToPlayerFunc = convertPayloadToPlayer

func fakeConvertPayloadToPlayerFunc(ctx context.Context, p PlayerEntityPayload) (*player.Player, errs.AppError) {
	return &player.Player{}, errs.ErrConvertingPayload
}

func restoreConvertPayloadToPlayerFunc(replace func(context.Context, PlayerEntityPayload) (*player.Player, errs.AppError)) {
	convertPayloadToPlayerFunc = replace
}

var convertAndValidatePayloadToTransferFunc = convertAndValidatePayloadToTransfer

func fakeConvertAndValidatePayloadToTransfer(ctx context.Context, t TransferEntityPayload) (*transfer.Transfer, errs.AppError) {
	return &transfer.Transfer{}, errs.ErrConvertingPayload
}

func restoreConvertAndValidatePayloadToTransfer(replace func(context.Context, TransferEntityPayload) (*transfer.Transfer, errs.AppError)) {
	convertAndValidatePayloadToTransferFunc = replace
}

var timeParse = time.Parse

func fakeTimeParse(layout string, value string) (time.Time, error) {
	return time.Time{}, errs.ErrParsingTime
}

func restoreTimeParse(replace func(layout string, value string) (time.Time, error)) {
	timeParse = replace
}

var convertPayloadToTournamentFunc = convertPayloadToTournament

func fakeConvertPayloadToTournament(ctx context.Context, t TournamentEntityPayload) (*tournament.Tournament, errs.AppError) {
	return &tournament.Tournament{}, errs.ErrConvertingPayload
}

func restoreConvertPayloadToTournament(replace func(ctx context.Context, t TournamentEntityPayload) (*tournament.Tournament, errs.AppError)) {
	convertPayloadToTournamentFunc = replace
}
