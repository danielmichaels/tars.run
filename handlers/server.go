package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net"
	"net/http"
	"shorty"
	"shorty/adapters"
	"strings"
	"time"
)

// respondWithJSON return json formatting and allow for custom response headers
func (s *apiServer) respondWithJSON(w http.ResponseWriter, i interface{}, status int) error {
	w.WriteHeader(status)
	e := ToJSON(i, w)
	return e
}

// getIP returns a users IP address. This is used to populate the DataPoints table.
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
	return "", fmt.Errorf("no valid ip found")
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

// Server interface implements a router type. Handles all HTTP requests to the server.
type Server interface {
	Router() *mux.Router
}

func (s apiServer) Router() *mux.Router {
	return s.router
}
