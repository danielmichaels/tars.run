package main

import (
	"bytes"
	"fmt"
	"github.com/danielmichaels/shortlink-go/cmd/api/config"
	"github.com/danielmichaels/shortlink-go/internal/data"
	"github.com/danielmichaels/shortlink-go/internal/logger"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestLinks(t *testing.T) {
	cfg := config.AppConfig()
	cfg.Db.DbName = "file::memory:?cache=shared"
	log := logger.New(os.Stdout, logger.LevelInfo)
	a := &application{
		config: cfg,
		logger: log,
		models: data.NewMockModels(),
	}

	t.Run("GET request to /v1/links should return 405 response", func(t *testing.T) {
		t.Helper()
		r, err := http.NewRequest("GET", "/v1/links", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		h := a.routes()
		h.ServeHTTP(rr, r)

		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("handler returned the wrong status code. got %v want %v", rr.Code, http.StatusMethodNotAllowed)
		}
	})
	t.Run("GET request to /v1/links/fake should return a 404", func(t *testing.T) {
		t.Helper()
		t.Log(cfg.Db)
		hash := "fake"
		r, err := http.NewRequest("GET", fmt.Sprintf("/v1/links/%s", hash), nil)
		if err != nil {
			t.Errorf("could not hit /v1/links/%s", hash)
		}
		rr := httptest.NewRecorder()

		h := a.routes()
		h.ServeHTTP(rr, r)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned the wrong status code. got %v want %v", rr.Code, http.StatusNotFound)
		}
	})
	t.Run("GET request to /v1/links/notfake should return a 200", func(t *testing.T) {
		t.Helper()
		hash := "notfake"
		r, err := http.NewRequest("GET", fmt.Sprintf("/v1/links/%s", hash), nil)
		if err != nil {
			t.Errorf("could not hit /v1/links/%s", hash)
		}
		rr := httptest.NewRecorder()

		h := a.routes()
		h.ServeHTTP(rr, r)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned the wrong status code. got %v want %v", rr.Code, http.StatusOK)
		}

		expected := `{"link":{"id":1,"created_at":"2020-01-01T00:00:00Z","original_url":"test.com","hash":"notfake"}}` + "\n"
		t.Log(rr.Body.String(), expected)
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body. got %v want %v", rr.Body.String(), expected)
		}
	})
	t.Run("POST a new link", func(t *testing.T) {
		t.Helper()
		link := `{"link":"test.com"}`
		r, err := http.NewRequest("POST", "/v1/links", bytes.NewBufferString(link))
		if err != nil {
			t.Error("could not hit /v1/links")
		}
		rr := httptest.NewRecorder()

		h := a.routes()
		h.ServeHTTP(rr, r)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned the wrong status code. got %v want %v", rr.Code, http.StatusOK)
		}
		// todo: FIX
		// Can't get the `hash` and `short_link` values to align as its dynamically created
		// on demand.
		//expected := `{"link":{"id":0,"created_at":"0001-01-01T00:00:00Z","original_url":"test.com","hash":"sZ1TY61fD5i"},"short_url":"http://localhost:1988/sZ1TY61fD5i"}` + "\n"
		//if rr.Body.String() != expected {
		//	t.Errorf("handler returned unexpected body. got %v want %v", rr.Body.String(), expected)
		//}
	})
}
