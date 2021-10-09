package handlers

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"net/http"
	"shorty/adapters"
	"strings"
	"time"
)

// NewLink accepts a post request of `link` and kicks off the encoding.
// returns: json object containing shortened link
func (s *apiServer) NewLink() http.HandlerFunc {
	type Request struct {
		Link string `json:"link"`
	}
	type Response struct {
		Data     adapters.Link `json:"data"`
		ShortUrl string        `json:"short_url"`
		Status   uint          `json:"status"`
		Msg      string        `json:"msg"`
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
		// todo: validate that link is a URL

		savedLink := adapters.Link{
			OriginalURL: *link,
			Hash:        hash,
		}

		resp := Response{
			Data:     savedLink,
			ShortUrl: savedLink.CreateShortLink(),
			Status:   http.StatusOK,
			Msg:      "successfully created link",
		}

		if err := s.db.Debug().Create(&savedLink).Error; err != nil {
			log.Fatalln("failed to create new short link in database", err)
		}
		var ex adapters.DataPoints
		s.db.Debug().First(&ex)
		log.Println(&ex)

		if err = s.respondWithJSON(w, &resp, http.StatusOK); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}

// AllLinks returns all links - testing and debug only todo
func (s *apiServer) AllLinks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := adapters.LinkModel{DB: s.db}
		links, err := l.All()
		if err != nil {
			log.Fatalln("failed to retrieve all Links from database")
			return
		}

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
	type Response struct {
		Result string `json:"result"`
		Msg    string `json:"msg"`
		Status uint   `json:"status"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var link adapters.Link
		hash := mux.Vars(r)["hash"]
		l := adapters.LinkModel{DB: s.db}
		link, err := l.Create(hash)
		//log.Println(err)
		//log.Println(link)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp := Response{
				Msg:    "domain not found",
				Status: http.StatusNotFound,
			}
			if err := s.respondWithJSON(w, &resp, http.StatusNotFound); err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			return
			//if err != nil {
			//	if errors.Is(err, gorm.ErrRecordNotFound) {
			//		resp := Response{
			//			Msg:    "domain not found",
			//			Status: http.StatusNotFound,
			//		}
			//		if err := s.respondWithJSON(w, &resp, http.StatusNotFound); err != nil {
			//			http.Error(w, "internal server error", http.StatusInternalServerError)
			//			return
			//		}
			//		return
			//	} else {
			//		http.Error(w, "failed to complete search", http.StatusInternalServerError)
			//		return
			//	}
		}
		domain := Response{
			Msg:    "success",
			Status: 200,
			Result: link.OriginalURL,
		}
		if !strings.HasPrefix(link.OriginalURL, "http://") && !strings.HasPrefix(link.OriginalURL, "https://") {
			domain.Result = "http://" + link.OriginalURL
		}

		// create a dataPoint each time its accessed
		ip, err := getIP(r)
		if err != nil {
			ip = "not found"
		}
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

		log.Println("domain resp:", domain)
		if err := s.respondWithJSON(w, domain, http.StatusOK); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}
}

// linkQuery looks up database for links matching a hash.
func (s *apiServer) linkQuery() http.HandlerFunc {
	type Response struct {
		Data   []adapters.DataPoints `json:"data"`
		Status uint                  `json:"status"`
		Msg    string                `json:"msg"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		hash := mux.Vars(r)["hash"]
		var data []adapters.DataPoints
		s.db.Debug().Where("link_id = ?", hash).Order("date_accessed desc").Find(&data)
		if len(data) < 1 {
			resp := Response{
				Status: 404,
				Msg:    "no data points",
				Data:   data,
			}
			if err := s.respondWithJSON(w, resp, http.StatusNotFound); err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			return
		}
		resp := Response{
			Data:   data,
			Status: 200,
			Msg:    "success",
		}

		if err := s.respondWithJSON(w, resp, http.StatusOK); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}
}
