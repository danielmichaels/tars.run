package handlers

import (
	"net/http"
	"shorty"
)

// Middleware that sets the `application/json` response type
func (s apiServer) jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware ensures that our frontend can talk to the server.
// In development, set by DEBUG=true, allow any origin.
func (s apiServer) CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get origin from config
		origin := shorty.AppConfig().Server.AllowedOrigins
		if shorty.AppConfig().Debug {
			// in dev accept any origin
			origin = r.Header.Get("Origin")
		}

		// Set headers
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS, DELETE, PATCH")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		// Next
		next.ServeHTTP(w, r)
	})
}
