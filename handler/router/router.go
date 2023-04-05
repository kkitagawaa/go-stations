package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/service"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()

	svc := service.NewTODOService(todoDB)

	mux.Handle("/healthz", handler.NewHealthzHandler())
	mux.Handle("/todos", handler.NewTODOHandler(svc))
	return mux
}
