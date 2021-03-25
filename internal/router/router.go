package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vimoppa/turl.to/internal/api"
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

	apiRoute.HandleFunc("/urls", api.AnyURLs(r.store)).Methods(http.MethodGet)
	apiRoute.HandleFunc("/urls", api.CreateURL(r.store)).Methods(http.MethodPost)
	apiRoute.HandleFunc("/urls/{hash}", api.FindOneURL(r.store)).Methods(http.MethodGet)

	// register subrouters.
	r.router.PathPrefix("/api").Handler(apiRoute)
}
