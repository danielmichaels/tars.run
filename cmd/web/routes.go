package main

import (
	"fmt"
	"github.com/danielmichaels/shortlink-go/internal/data"
	"github.com/danielmichaels/shortlink-go/internal/validator"
	"github.com/danielmichaels/shortlink-go/ui"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
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
	fileServer := http.FileServer(neuteredFileSystem{http.FS(ui.Files)})
	r.Handle("/static/*", fileServer)

	// Routes
	r.Get("/", app.handleHomepage())
	//r.Get("/v1/healthcheck", app.healthcheckHandler)
	//r.Get("/v1/links/{hash}", app.showLinkHandler())
	r.Post("/v1/links", app.handleCreateLink())
	//r.Get("/v1/links/{hash}/analytics", app.showLinkAnalyticsHandler())

	return r
}

func (app *application) handleHomepage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		links, err := app.models.Links.Get("yUitgPWMVyg")
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		app.logger.Info().Msgf("%#v", links)
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
