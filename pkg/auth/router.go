package auth

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/rafaelsanzio/go-flashscore/pkg/auth/handlers"
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
	{Name: "Health Check Auth API", Methods: []string{http.MethodGet}, Path: "/auth-ok", Handler: handlers.HandleAdapter(handlers.HandleAuthOK)},

	{Name: "Creating API Key", Methods: []string{http.MethodPost}, Path: "/apikey", Handler: handlers.HandleAdapter(handlers.HandleCreateAPIKey)},
}
