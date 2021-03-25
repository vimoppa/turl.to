package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vimoppa/turl.to/internal/storage"
)

// Router contains http mux and storage connector
type Router struct {
	router *mux.Router
	store  storage.Accessor
}

// New creates a new router
func New(s storage.Accessor) *Router {
	r := &Router{
		router: mux.NewRouter(),
		store:  s,
	}
	r.initRoutes()
	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func (r *Router) initRoutes() {
	// api route group
	apiRoute := mux.NewRouter().PathPrefix("/api").Subrouter()

	// register subrouters.
	r.router.PathPrefix("/api").Handler(apiRoute)
}
