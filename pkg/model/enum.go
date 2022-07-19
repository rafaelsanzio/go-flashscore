package model

import (
	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
)

const (
	KafkaTopicTransfer    = "transfer"
	KafkaTopicMatchEvents = "match-events"
)

var Topics = []string{KafkaTopicTransfer, KafkaTopicMatchEvents}

type ActionType string

var (
	actionTypes = make(map[string]ActionType, 5)
)

func actionType(name string) ActionType {
	i := ActionType(name)
	actionTypes[name] = i
	return i
}

func (i *ActionType) UnmarshalText(data []byte) error {

	str := string(data)

	val, ok := actionTypes[str]

	if ok {
		*i = val
		return nil
	}

	return errs.ErrInvalidActionType.Throwf(applog.Log, "type: %s", str)
}

var (
	ActionUnknown = actionType("Unknown")

	ActionUpdateTeamPlayer = actionType("UpdateTeamPlayer")
	ActionGameEvents       = actionType("ActionGameEvents")
)

type EventsMatchType string

var (
	eventsMatchTypes = make(map[string]EventsMatchType, 7)
)

func eventsMatchType(name string) EventsMatchType {
	i := EventsMatchType(name)
	eventsMatchTypes[name] = i
	return i
}

func (i *EventsMatchType) UnmarshalText(data []byte) error {

	str := string(data)

	val, ok := eventsMatchTypes[str]

	if ok {
		*i = val
		return nil
	}

	return errs.ErrInvalidActionType.Throwf(applog.Log, "type: %s", str)
}

var (
	EventStart        = eventsMatchType("Start")
	EventGoal         = eventsMatchType("Goal")
	EventHalftime     = eventsMatchType("Halftime")
	EventExtratime    = eventsMatchType("Extratime")
	EventSubstitution = eventsMatchType("Substitution")
	EventWarning      = eventsMatchType("Warning")
	EventFinish       = eventsMatchType("Finish")
)

type MatchStatus string

var (
	matchStatusTypes = make(map[string]MatchStatus, 4)
)

func matchStatusType(name string) MatchStatus {
	i := MatchStatus(name)
	matchStatusTypes[name] = i
	return i
}

func (i *MatchStatus) UnmarshalText(data []byte) error {

	str := string(data)

	val, ok := matchStatusTypes[str]

	if ok {
		*i = val
		return nil
	}

	return errs.ErrInvalidActionType.Throwf(applog.Log, "type: %s", str)
}

var (
	MatchStatusNotStart   = matchStatusType("NotStarted")
	MatchStatusInProgress = matchStatusType("InProgress")
	MatchStatusHalftime   = matchStatusType("MatchHalftime")
	MatchStatusFinished   = matchStatusType("Finished")
)

type Warnings string

var (
	warningsTypes = make(map[string]Warnings, 2)
)

func warningsType(name string) Warnings {
	i := Warnings(name)
	warningsTypes[name] = i
	return i
}

func (i *Warnings) UnmarshalText(data []byte) error {

	str := string(data)

	val, ok := warningsTypes[str]

	if ok {
		*i = val
		return nil
	}

	return errs.ErrInvalidActionType.Throwf(applog.Log, "type: %s", str)
}

var (
	WarningRedCard    = warningsType("RedCard")
	WarningYellowCard = warningsType("YellowCard")
)
