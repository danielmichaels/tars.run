package handlers

import "net/http"

// All routes for the server are found here.
func (s *apiServer) routes() {
	api := s.router.PathPrefix("/api").Subrouter()
	api.Use(s.jsonContentTypeMiddleware)

	// for testing, this should only return a users links in future
	api.HandleFunc("/links/all", s.AllLinks()).Methods(http.MethodGet)
	api.HandleFunc("/links/analytics/all", s.AllDataPoints()).Methods(http.MethodGet)

	// create new links
	api.HandleFunc("/new", s.NewLink()).Methods(http.MethodPost)

	api.HandleFunc("/{hash}/query", s.linkQuery()).Methods(http.MethodGet)

	api.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
	})

	core := s.router.PathPrefix("/").Subrouter()
	// redirect matching shortened urls
	core.HandleFunc("/nothing-found", s.nothingFound()).Methods(http.MethodGet)
	core.HandleFunc("/{hash}", s.resolveHash()).Methods(http.MethodGet)
}
