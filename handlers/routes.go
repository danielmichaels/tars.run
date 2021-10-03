package handlers

import (
	"net/http"
)

// All routes for the server are found here.
func (s *apiServer) routes() {
	api := s.router.PathPrefix("/api").Subrouter()
	api.Use(s.CORSMiddleware)
	api.Use(s.jsonContentTypeMiddleware)

	// for testing, this should only return a users links in future
	api.HandleFunc("/links/all", s.AllLinks()).Methods(http.MethodGet)
	api.HandleFunc("/links/analytics/all", s.AllDataPoints()).Methods(http.MethodGet)

	// create new links
	//api.HandleFunc("/new", allowOptions()).Methods(http.MethodOptions)
	api.HandleFunc("/new", s.NewLink()).Methods(http.MethodPost, http.MethodOptions)

	api.HandleFunc("/{hash}/query", s.linkQuery()).Methods(http.MethodGet)
	api.HandleFunc("/{hash}", s.resolveHash()).Methods(http.MethodGet)

	api.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
	})

	core := s.router.PathPrefix("/").Subrouter()
	core.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
	})
}
