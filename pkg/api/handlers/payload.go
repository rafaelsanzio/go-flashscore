package handlers

import (
	"github.com/rafaelsanzio/go-flashscore/pkg/model"
	"github.com/rafaelsanzio/go-flashscore/pkg/money"
)

type TeamEntityPayload struct {
	Name      string `json:"name"`
	ShortCode string `json:"short_code"`
	Country   string `json:"country"`
	City      string `json:"city"`
}

type PlayerEntityPayload struct {
	Name         string `json:"name"`
	Team         string `json:"team"`
	Country      string `json:"country"`
	BirthdayDate string `json:"birthday_date"`
}

type TransferEntityPayload struct {
	Player         string      `json:"player"`
	TeamOrigin     string      `json:"team_origin"`
	TeamDestiny    string      `json:"team_destiny"`
	DateOfTransfer string      `json:"date_of_transfer"`
	Amount         money.Money `json:"amount"`
}

type TournamentEntityPayload struct {
	Name  string   `json:"name"`
	Teams []string `json:"teams"`
}

type AddTeamsTournamentEntityPayload struct {
	Teams []string `json:"teams"`
}

type MatchEntityPayload struct {
	HomeTeam    string `json:"home_team"`
	AwayTeam    string `json:"away_team"`
	DateOfMatch string `json:"date_of_match"`
	TimeOfMatch string `json:"time_of_match"`
}

type MatchGoalEntityPayload struct {
	TeamScore string `json:"team_score"`
	Player    string `json:"player"`
	Minute    int    `json:"minute"`
}

type ExtraTimeEntityPayload struct {
	Extratime int `json:"extratime"`
}

type MatchSubstitutionPayload struct {
	Team      string `json:"team"`
	PlayerOut string `json:"player_out"`
	PlayerIn  string `json:"player_in"`
	Minute    int    `json:"minute"`
}

type MatchWarningPayload struct {
	Team    string         `json:"team"`
	Player  string         `json:"player"`
	Warning model.Warnings `json:"warning"`
	Minute  int            `json:"minute"`
}
