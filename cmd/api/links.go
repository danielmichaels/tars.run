package main

import (
	"errors"
	"fmt"
	"github.com/danielmichaels/shortlink-go/internal/data"
	"github.com/danielmichaels/shortlink-go/internal/validator"
	"github.com/gorilla/mux"
	"github.com/tomasen/realip"
	"net/http"
)

func (app *application) showLinkHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := mux.Vars(r)["hash"] // todo: validator

		link, err := app.models.Links.Get(hash)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		analytic := data.Analytics{
			Ip:        realip.FromRequest(r),
			UserAgent: r.UserAgent(),
			LinkID:    uint64(link.ID),
		}

		// this constitutes a query on the link, so we save this to the Analytics table.
		err = app.models.Analytics.Insert(&analytic)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		err = app.writeJSON(w, http.StatusOK, envelope{"link": link}, nil)
		app.logger.PrintInfo("link data", map[string]string{
			"link.hash":         link.Hash,
			"link.original_url": link.OriginalURL,
		})
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}

}

func (app *application) createLinkHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Link string `json:"link"`
		}

		err := app.readJSON(w, r, &input)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		link := &data.Link{
			OriginalURL: input.Link,
			Hash:        data.CreateURL(),
		}

		// todo Validate

		err = app.models.Links.Insert(link)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		headers := make(http.Header)
		headers.Set("Location", fmt.Sprintf("/v1/links/%s", link.Hash))

		err = app.writeJSON(w, http.StatusCreated, envelope{"link": link, "short_url": link.CreateShortLink()}, headers)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}
}

func (app *application) showLinkAnalyticsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := mux.Vars(r)["hash"]
		var input struct {
			data.Filters
		}
		v := validator.New()

		qs := r.URL.Query()

		input.Filters.Page = app.readInt(qs, "page", 1, v)
		input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
		input.Filters.Sort = app.readString(qs, "sort", "date_accessed")
		input.Filters.SortSafeList = []string{"id", "date_accessed", "user_agent", "ip", "-id", "-date_accessed", "-user_agent", "-ip"}

		if data.ValidateFilters(v, input.Filters); !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		links, metadata, err := app.models.Analytics.GetAllForLink(hash, input.Filters)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		err = app.writeJSON(w, http.StatusOK, envelope{"metadata": metadata, "analytics": links}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}
}
