package handlers

import (
	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"shorty/adapters"
	"testing"
)

var mockLinks = []adapters.Link{
	{OriginalURL: "https://test.com", Hash: "12345678", Data: []adapters.DataPoints{}},
	{OriginalURL: "mudmap.io", Hash: "abcdefgh", Data: []adapters.DataPoints{}},
}

func MockDBConnection() *gorm.DB {
	log.Println("Connecting to mock database")
	database := "file::memory:?cache=shared"
	//database := "test_mock.db"
	db, err := gorm.Open(sqlite.Open(database), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to adapters")
	}
	adapters.InitialMigrations(database)
	seedDatabase(database)
	return db
}

func seedDatabase(database string) {
	db, err := gorm.Open(sqlite.Open(database), &gorm.Config{})
	if err := db.Migrator().DropTable(adapters.Link{}, adapters.DataPoints{}); err != nil {
		log.Fatalln("could not drop tables in `seedDatabase` func")
	}
	adapters.InitialMigrations(database)
	if err != nil {
		log.Fatalln("failed to connect to mock database during seed")
	}
	if err := db.Debug().Create(&mockLinks).Error; err != nil {
		log.Fatalln("failed to create `mockLinks` during seed")
	}
}

func TestServer_Router(t *testing.T) {
	// set up the server

	router := mux.NewRouter()
	s := apiServer{
		router:       router,
		ReadTimeout:  1,
		WriteTimeout: 5,
		IdleTimeout:  15,
		db:           MockDBConnection(),
	}
	s.routes()
	//MockDBConnection()

	t.Run("API returns invalid path if no router", func(t *testing.T) {
		t.Helper()
		request, err := http.NewRequest(http.MethodGet, "/api/", nil)
		if err != nil {
			t.Fatal(err)
		}
		resp := httptest.NewRecorder()
		s.router.ServeHTTP(resp, request)
		checkResponseCode(t, resp.Code, http.StatusBadRequest)
	})
	t.Run("Core returns invalid path if no router", func(t *testing.T) {
		t.Helper()
		request, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		resp := httptest.NewRecorder()
		s.router.ServeHTTP(resp, request)
		checkResponseCode(t, resp.Code, http.StatusBadRequest)
	})
	t.Run("GET all links", func(t *testing.T) {
		t.Helper()
		request, err := http.NewRequest(http.MethodGet, "/api/links/all", nil)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		resp := httptest.NewRecorder()
		s.router.ServeHTTP(resp, request)
		checkResponseCode(t, resp.Code, http.StatusOK)
		expected := len(mockLinks)
		var actual []adapters.Link
		err = FromJSON(&actual, resp.Body)
		if err != nil {
			t.Fatalf("Unable to parse response from server %q into slice of adapters.Link, '%v'", resp.Body, err)
		}
		t.Log(resp.Body.String())
		t.Logf("%#v", actual)
		if len(actual) != expected {
			t.Errorf("actual did not return correct number of rows")
		}
	})
	t.Run("resolveHash returns a single link", func(t *testing.T) {
		t.Helper()

		request, err := http.NewRequest(http.MethodGet, "/api/12345678", nil)
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		resp := httptest.NewRecorder()
		s.router.ServeHTTP(resp, request)
		checkResponseCode(t, resp.Code, http.StatusOK)
		//expected := "test.com"
		expected := mockLinks[0].OriginalURL

		t.Log("expected", expected) // DEBUG

		var actual adapters.Link
		err = FromJSON(&actual, resp.Body)
		t.Logf("%#v", resp)
		t.Logf("%#v", resp.Body.String())
		t.Log("actual.URL", actual.OriginalURL)
		t.Log("actual", actual)
		t.Log(actual.OriginalURL)
		t.Log("len", len(actual.OriginalURL))
		if err != nil {
			t.Fatalf("Unable to parse response. Expected %q, actual %v", resp.Body, err)
		}
		got := actual.Hash
		if expected != got {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		//if !reflect.DeepEqual(actual, expected) {
		//	t.Errorf("expected %v, actual %v", expected, actual)
		//}
	})
	t.Run("resolveHash fails to return a single link and gets 404", func(t *testing.T) {
		t.Helper()

		request, err := http.NewRequest(http.MethodGet, "/api/doesnotexist", nil)
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		resp := httptest.NewRecorder()
		s.router.ServeHTTP(resp, request)
		checkResponseCode(t, resp.Code, http.StatusNotFound)

	})
}

// checkResponseCode testing utility for asserting the status code
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}

}
