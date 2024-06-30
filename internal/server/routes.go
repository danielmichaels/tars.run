package server

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/danielmichaels/shortlink-go/assets/view/pages"
	"github.com/danielmichaels/shortlink-go/internal/templates"

	"github.com/danielmichaels/shortlink-go/assets"
	"github.com/danielmichaels/shortlink-go/internal/data"
	"github.com/danielmichaels/shortlink-go/internal/validator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/tomasen/realip"
)

func (app *Application) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.RealIP)
	r.Use(middleware.Compress(5))
	r.Use(httplog.RequestLogger(app.Logger, []string{
		"/healthz",
		"/static/img/logo.png",
		"/static/icons/favicon.ico",
		"/static/css/theme.css",
		"/static/js/bundle.js",
		"/static/js/local-storage.js",
		"/static/js/htmx.min.js",
	}))
	r.Use(middleware.Heartbeat("/healthz"))

	r.NotFound(app.notFound)
	r.MethodNotAllowed(app.methodNotAllowed)

	fileServer := http.FileServer(http.FS(assets.Files))
	r.Handle("/static/*", fileServer)

	r.Post("/v1/links", app.handleCreateLink())

	r.Get("/", app.handleHome)
	r.Get("/{hash}", app.handleRedirect)
	r.Get("/{hash}/analytics", app.handleAnalytics)
	return r
}

func (app *Application) handleHome(w http.ResponseWriter, r *http.Request) {
	err := templates.Render(r.Context(), w, 200, pages.Home())
	if err != nil {
		return
	}
}

func (app *Application) handleRedirect(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	link, err := app.Models.Links.Get(hash)
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
	err = app.Models.Analytics.Insert(&analytic)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.Logger.Info().Msgf("redirect: %s-%s", link.Hash, link.OriginalURL)
	// Use a temporary redirect status in case we want to support changing
	// redirect targets in the future.
	http.Redirect(w, r, link.OriginalURL, http.StatusTemporaryRedirect)
	app.Logger.Info().Msgf("redirect: %s-%s", link.Hash, link.OriginalURL)
}
func (app *Application) handleAnalytics(w http.ResponseWriter, r *http.Request) {
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
	input.Filters.SortSafeList = []string{
		"id",
		"date_accessed",
		"user_agent",
		"ip",
		"-id",
		"-date_accessed",
		"-user_agent",
		"-ip",
	}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	analytics, metadata, err := app.Models.Analytics.GetAllForLink(hash, input.Filters)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFound(w, r)
		default:
			app.serverError(w, r, err)
		}
		return
	}
	// Return newest first
	slices.Reverse(analytics)

	link, err := app.Models.Links.Get(hash)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	fmt.Printf("%v (%T)\n", link.OriginalURL, link.OriginalURL)
	a := pages.AnalyticsData{
		Title:     "Analytics",
		Link:      link,
		Analytics: analytics,
		Metadata:  metadata,
	}

	_ = templates.Render(r.Context(), w, 200, pages.Analytics(a))
}

func (app *Application) notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	_, _ = w.Write([]byte("404 Not Found"))

}
func (app *Application) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(405)
	_, _ = w.Write([]byte("Method Not Allowed"))
}

func (app *Application) handleCreateLink() http.HandlerFunc {
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
		err = app.Models.Links.Insert(link)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		headers := make(http.Header)
		headers.Set("Location", fmt.Sprintf("/v1/links/%s", link.Hash))

		err = app.writeJSON(
			w,
			http.StatusCreated,
			M{"link": link, "short_url": link.CreateShortLink()},
			headers,
		)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}
}
