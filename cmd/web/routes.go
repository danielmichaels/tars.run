package main

import (
	"errors"
	"fmt"
	"github.com/danielmichaels/shortlink-go/assets"
	"github.com/danielmichaels/shortlink-go/internal/data"
	"github.com/danielmichaels/shortlink-go/internal/validator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/tomasen/realip"
	"net/http"
	"time"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	// Middleware
	r.Use(middleware.RealIP)
	r.Use(httplog.RequestLogger(app.logger))
	r.Use(app.recoverPanic)
	r.NotFound(app.notFound)
	r.MethodNotAllowed(app.methodNotAllowed)
	// fileServer for static assets
	fileServer := http.FileServer(neuteredFileSystem{http.FS(assets.Files)})
	r.Handle("/static/*", fileServer)

	// Routes
	r.Get("/", app.handleHomepage())
	r.Get("/{hash}", app.handleRedirectLink())
	r.Get("/{hash}/analytics", app.handleLinkAnalytics())

	r.Post("/v1/links", app.handleCreateLink())

	return r
}
func (app *application) notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte("404 Not Found"))

}
func (app *application) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(405)
	w.Write([]byte("Method Not Allowed"))
}
func (app *application) handleHomepage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.render(w, r, "home.page.tmpl", &templateData{
			Title: "Home",
		})
	}
}

func (app *application) handleCreateLink() http.HandlerFunc {
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

		v := validator.New()

		if data.ValidateURL(v, input.Link); !v.Valid() {
			// prepend http:// to the link provided by the user.
			link.OriginalURL = fmt.Sprintf("http://%s", input.Link)
		}

		time.Sleep(1 * time.Second)
		err = app.models.Links.Insert(link)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		headers := make(http.Header)
		headers.Set("Location", fmt.Sprintf("/v1/links/%s", link.Hash))

		err = app.writeJSON(w, http.StatusCreated, envelope{"link": link, "short_url": link.CreateShortLink()}, headers)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}
}

func (app *application) handleRedirectLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := chi.URLParam(r, "hash")

		link, err := app.models.Links.Get(hash)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				// record not found should not display server response which is unhelpful
				// it instead redirects to the /404.
				app.notFound(w, r)
			default:
				app.serverError(w, r, err)
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
			app.serverError(w, r, err)
			return
		}
		app.logger.Info().Msgf("redirect: %s-%s", link.Hash, link.OriginalURL)
		// Use a temporary redirect status in case we want to support changing
		// redirect targets in the future.
		http.Redirect(w, r, link.OriginalURL, http.StatusTemporaryRedirect)
		app.logger.Info().Msgf("redirect: %s-%s", link.Hash, link.OriginalURL)
	}
}

func (app *application) handleLinkAnalytics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := chi.URLParam(r, "hash")
		var input struct {
			data.Filters
		}
		v := validator.New()

		qs := r.URL.Query()

		// Allow ability to paginate response in the future.
		input.Filters.Page = validator.ReadInt(qs, "page", 1, v)
		input.Filters.PageSize = validator.ReadInt(qs, "page_size", 20, v)
		input.Filters.Sort = validator.ReadString(qs, "sort", "date_accessed")
		input.Filters.SortSafeList = []string{"id", "date_accessed", "user_agent", "ip", "-id", "-date_accessed", "-user_agent", "-ip"}

		if data.ValidateFilters(v, input.Filters); !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		analytics, metadata, err := app.models.Analytics.GetAllForLink(hash, input.Filters)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFound(w, r)
			default:
				app.serverError(w, r, err)
			}
			return
		}
		link, err := app.models.Links.Get(hash)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		app.logger.Info().Msgf("%#v", analytics)
		app.render(w, r, "analytics.page.tmpl", &templateData{
			Title:     "Analytics",
			Link:      link,
			Analytics: analytics,
			Metadata:  metadata,
		})
	}
}
