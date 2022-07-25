package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/rafaelsanzio/go-flashscore/pkg/api/handlers"
)

func NewRouter() http.Handler {
	router := mux.NewRouter()

	for _, r := range routes {
		router.Methods(r.Methods...).Path(r.Path).Name(r.Name).HandlerFunc(r.Handler)
	}

	return router
}

type Route struct {
	Name    string
	Methods []string
	Path    string
	Handler http.HandlerFunc
}

var routes = []Route{
	{Name: "Health Check API", Methods: []string{http.MethodGet}, Path: "/ok", Handler: handlers.HandleAPIOK},

	// Team
	{Name: "Creating a team", Methods: []string{http.MethodPost}, Path: "/teams", Handler: handlers.HandleAdapter(handlers.HandlePostTeam)},
	{Name: "Listing all teams", Methods: []string{http.MethodGet}, Path: "/teams", Handler: handlers.HandleAdapter(handlers.HandleListTeam)},
	{Name: "Getting a team", Methods: []string{http.MethodGet}, Path: "/teams/{id}", Handler: handlers.HandleAdapter(handlers.HandleGetTeam)},
	{Name: "Updating a team", Methods: []string{http.MethodPut}, Path: "/teams/{id}", Handler: handlers.HandleAdapter(handlers.HandleUpdateTeam)},
	{Name: "Deleting a team", Methods: []string{http.MethodDelete}, Path: "/teams/{id}", Handler: handlers.HandleAdapter(handlers.HandleDeleteTeam)},

	// Player
	{Name: "Creating a player", Methods: []string{http.MethodPost}, Path: "/players", Handler: handlers.HandleAdapter(handlers.HandlePostPlayer)},
	{Name: "Listing all players", Methods: []string{http.MethodGet}, Path: "/players", Handler: handlers.HandleAdapter(handlers.HandleListPlayer)},
	{Name: "Getting a player", Methods: []string{http.MethodGet}, Path: "/players/{id}", Handler: handlers.HandleAdapter(handlers.HandleGetPlayer)},
	{Name: "Updating a player", Methods: []string{http.MethodPut}, Path: "/players/{id}", Handler: handlers.HandleAdapter(handlers.HandleUpdatePlayer)},
	{Name: "Deleting a player", Methods: []string{http.MethodDelete}, Path: "/players/{id}", Handler: handlers.HandleAdapter(handlers.HandleDeletePlayer)},

	// Transfer
	{Name: "Creating a transfer", Methods: []string{http.MethodPost}, Path: "/transfers", Handler: handlers.HandleAdapter(handlers.HandlePostTransfer)},
	{Name: "Getting a transfer", Methods: []string{http.MethodGet}, Path: "/transfers/{id}", Handler: handlers.HandleAdapter(handlers.HandleGetTransfer)},
	{Name: "Listing all transfers", Methods: []string{http.MethodGet}, Path: "/transfers", Handler: handlers.HandleAdapter(handlers.HandleListTransfer)},

	// Tournament
	{Name: "Creating a tournament", Methods: []string{http.MethodPost}, Path: "/tournaments", Handler: handlers.HandleAdapter(handlers.HandlePostTournament)},
	{Name: "Listing all tournaments", Methods: []string{http.MethodGet}, Path: "/tournaments", Handler: handlers.HandleAdapter(handlers.HandleListTournament)},
	{Name: "Getting a tournament", Methods: []string{http.MethodGet}, Path: "/tournaments/{id}", Handler: handlers.HandleAdapter(handlers.HandleGetTournament)},
	{Name: "Updating a tournament", Methods: []string{http.MethodPut}, Path: "/tournaments/{id}", Handler: handlers.HandleAdapter(handlers.HandleUpdateTournament)},
	{Name: "Deleting a tournament", Methods: []string{http.MethodDelete}, Path: "/tournaments/{id}", Handler: handlers.HandleAdapter(handlers.HandleDeleteTournament)},

	// Tournament -> Teams
	{Name: "Adding teams to a tournament", Methods: []string{http.MethodPost}, Path: "/tournaments/{id}/add-teams", Handler: handlers.HandleAdapter(handlers.HandleAddTeamsTournament)},

	// Tournament -> Matches
	{Name: "Creating a match to a tournament", Methods: []string{http.MethodPost}, Path: "/tournaments/{id}/matches", Handler: handlers.HandleAdapter(handlers.HandlePostMatch)},
	{Name: "Listing all match from tournament", Methods: []string{http.MethodGet}, Path: "/tournaments/{id}/matches", Handler: handlers.HandleAdapter(handlers.HandleListMatch)},
	{Name: "Getting a match from tournament", Methods: []string{http.MethodGet}, Path: "/tournaments/{id}/matches/{match_id}", Handler: handlers.HandleAdapter(handlers.HandleGetMatch)},
	{Name: "Deleting a match from tournament", Methods: []string{http.MethodDelete}, Path: "/tournaments/{id}/matches/{match_id}", Handler: handlers.HandleAdapter(handlers.HandleDeleteMatch)},

	// Tournament -> Matches -> Events
	{Name: "Creating an event to start a match", Methods: []string{http.MethodPost}, Path: "/tournaments/{id}/matches/{match_id}/events/start", Handler: handlers.HandleAdapter(handlers.HandlePostMatchStart)},
	{Name: "Creating an event to score a goal in a match", Methods: []string{http.MethodPost}, Path: "/tournaments/{id}/matches/{match_id}/events/goal", Handler: handlers.HandleAdapter(handlers.HandlePostMatchGoal)},
	{Name: "Creating an event to halftime a match", Methods: []string{http.MethodPost}, Path: "/tournaments/{id}/matches/{match_id}/events/halftime", Handler: handlers.HandleAdapter(handlers.HandlePostMatchHalftime)},
	{Name: "Creating an event to substitution players in a match", Methods: []string{http.MethodPost}, Path: "/tournaments/{id}/matches/{match_id}/events/substitution", Handler: handlers.HandleAdapter(handlers.HandlePostMatchSubstitution)},
	{Name: "Creating an event to add a warning in a match", Methods: []string{http.MethodPost}, Path: "/tournaments/{id}/matches/{match_id}/events/warning", Handler: handlers.HandleAdapter(handlers.HandlePostMatchWarning)},
	{Name: "Creating an event to add extratime in a match", Methods: []string{http.MethodPost}, Path: "/tournaments/{id}/matches/{match_id}/events/extratime", Handler: handlers.HandleAdapter(handlers.HandlePostMatchExtratime)},
	{Name: "Creating an event to finish a match", Methods: []string{http.MethodPost}, Path: "/tournaments/{id}/matches/{match_id}/events/finish", Handler: handlers.HandleAdapter(handlers.HandlePostMatchFinish)},
}
