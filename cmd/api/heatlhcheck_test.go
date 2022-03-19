package main

import (
	"fmt"
	"github.com/danielmichaels/shortlink-go/internal/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthcheckHandler(t *testing.T) {
	cfg := config.AppConfig()
	cfg.Debug = false
	app := &application{config: cfg}
	t.Run("get request to endpoint successful", func(t *testing.T) {
		t.Helper()
		r, err := http.NewRequest("GET", "/v1/healthcheck", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		h := app.routes()
		h.ServeHTTP(rr, r)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned the wrong status code: got %v want %v", status, http.StatusOK)
		}
		// adding the '\n' to writeJSON means we have to include it here or it's not a match
		expected := `{"status":"available","system_info":{"debug":"false","version":"1.0.0"}}` + "\n"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})
	t.Run("post request to healthcheck returns method not allowed", func(t *testing.T) {
		t.Helper()
		r, err := http.NewRequest(http.MethodPost, "/v1/healthcheck", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		h := app.routes()
		h.ServeHTTP(rr, r)

		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("handler expected to return a 405 Method Not Allowed but didnt. got %v want %v", status, http.StatusMethodNotAllowed)
		}

		expected := fmt.Sprintf(`{"error":"the method %s is not supported for this resource"}`, r.Method) + "\n"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body. got %v want %v", rr.Body.String(), expected)
		}
	})
}
