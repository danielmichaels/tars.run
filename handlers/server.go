package handlers

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"net"
	"net/http"
	"shorty"
	"shorty/adapters"
	"strings"
	"time"
)

// Server interface implements a router type. Handles all HTTP requests to the server.
type Server interface {
	Router() *mux.Router
}

func (s apiServer) Router() *mux.Router {
	return s.router
}

// respondWithJSON return json formatting and allow for custom response headers
func (s *apiServer) respondWithJSON(w http.ResponseWriter, i interface{}, status int) error {
	w.WriteHeader(status)
	e := ToJSON(i, w)
	return e
}

// NewLink accepts a post request of `link` and kicks off the encoding.
// returns: json object containing shortened link
func (s *apiServer) NewLink() http.HandlerFunc {
	type Request struct {
		Link string `json:"link"`
	}
	var request Request
	return func(w http.ResponseWriter, r *http.Request) {
		err := FromJSON(&request, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		link := &request.Link
		hash := adapters.CreateURL()

		response := adapters.Link{
			OriginalURL: *link,
			Hash:        hash,
		}

		if err := s.db.Debug().Create(&response).Error; err != nil {
			log.Fatalln("failed to create new short link in database", err)
		}
		var ex adapters.DataPoints
		s.db.Debug().First(&ex)
		log.Println(&ex)

		err = s.respondWithJSON(w, &response, http.StatusOK)
	}
}

// AllLinks returns all links - testing and debug only todo
func (s *apiServer) AllLinks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var links []adapters.Link
		s.db.Find(&links)
		log.Println(links)

		if err := s.respondWithJSON(w, links, http.StatusOK); err != nil {
			http.Error(w, "error occurred", http.StatusInternalServerError)
			return
		}
	}
}

// AllDataPoints returns all data points for links - testing and debug only todo
func (s *apiServer) AllDataPoints() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data []adapters.DataPoints
		s.db.Debug().Find(&data)
		fmt.Println(data)

		if err := s.respondWithJSON(w, data, http.StatusOK); err != nil {
			http.Error(w, "error occurred", http.StatusInternalServerError)
			return
		}
	}
}

func (s *apiServer) resolveHash() http.HandlerFunc {
	var link adapters.Link
	return func(w http.ResponseWriter, r *http.Request) {
		hash := mux.Vars(r)["hash"]
		// todo: its looking up any asset i.e. favicon/css etc - should only lookup the {hash}
		log.Println("hash", hash)
		err := s.db.Debug().Where("hash = ?", hash).First(&link).Error
		if err != nil {
			log.Println("NOT FOUND")
			if errors.Is(err, gorm.ErrRecordNotFound) {
				domain := fmt.Sprintf("%s/nothing-found", shorty.AppConfig().Server.Domain)
				http.Redirect(w, r, domain, http.StatusFound)
				return
			} else {
				http.Error(w, "failed to complete search", http.StatusInternalServerError)
				return
			}
		}
		log.Println(link.Hash)
		// http://localhost:1987/a7GCBojygRk // mudmap.io
		// http://localhost:1987/Boo9ohcsEsr // danielms
		var domain string
		if !strings.HasPrefix(link.OriginalURL, "http://") && !strings.HasPrefix(link.OriginalURL, "https://") {
			domain = "http://" + link.OriginalURL
			log.Println("prepending domain", domain)
		} else {
			domain = link.OriginalURL
			log.Println("domain has prefix already", domain)
		}

		w.Header().Add("Content-Type", "text/plain")
		w.Header().Set("Location", domain)

		// create a dataPoint each time its accessed
		ip, err := getIP(r)
		if err != nil {
			ip = "not found"
		}
		fmt.Println(ip)
		data := &adapters.DataPoints{
			LinkID:       link.Hash,
			DateAccessed: time.Now().UTC(),
			Ip:           ip,
			UserAgent:    r.UserAgent(),
			Location:     r.Referer(), // this is wrong, needs an IP lookup
		}
		if err := s.db.Debug().Create(&data).Error; err != nil {
			log.Println("failed to save datapoint", err)
		}

		http.Redirect(w, r, domain, http.StatusSeeOther)
	}
}

func (s *apiServer) linkQuery() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := mux.Vars(r)["hash"]
		log.Println("query", hash)
		var data []adapters.DataPoints
		s.db.Debug().Where("link_id = ?", hash).Order("date_accessed desc").Find(&data)
		fmt.Println(data)

		if err := s.respondWithJSON(w, data, http.StatusOK); err != nil {
			http.Error(w, "error occurred", http.StatusInternalServerError)
			return
		}
	}
}

func (s *apiServer) nothingFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("nothing found"))

	}
}

func getIP(r *http.Request) (string, error) {
	//Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	//Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	return "", fmt.Errorf("No valid ip found")
}

func NewServer() Server {
	s := apiServer{
		router:       mux.NewRouter(),
		ReadTimeout:  shorty.AppConfig().Server.TimeoutRead,
		WriteTimeout: shorty.AppConfig().Server.TimeoutWrite,
		IdleTimeout:  shorty.AppConfig().Server.TimeoutIdle,
		db:           adapters.DBConnection(),
	}
	s.routes()
	return s
}

type apiServer struct {
	router       *mux.Router
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	db           *gorm.DB
}
