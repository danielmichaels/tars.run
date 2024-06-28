package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/danielmichaels/shortlink-go/internal/config"
	"github.com/danielmichaels/shortlink-go/internal/data"
	"github.com/rs/zerolog"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Application struct {
	Config *config.Conf
	Logger zerolog.Logger
	//DB     *store.Queries
	Models   data.Models
	Template map[string]*template.Template
}

type templateData struct {
	Title     string
	AppUrl    string
	Names     userTemplateData
	Link      *data.Link
	Analytics []*data.Analytics
	Metadata  data.Metadata
}

type userTemplateData struct {
	AppName   string
	Github    string
	Twitter   string
	Plausible string
}

func (app *Application) Serve(ctx context.Context) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.Server.Port),
		Handler:      app.routes(),
		IdleTimeout:  app.Config.Server.TimeoutIdle,
		ReadTimeout:  app.Config.Server.TimeoutRead,
		WriteTimeout: app.Config.Server.TimeoutWrite,
	}
	app.Logger.Info().Msgf("HTTP webserver listening on '%d'", app.Config.Server.Port)
	wg := sync.WaitGroup{}
	shutdownError := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.Logger.Warn().Str("signal", s.String()).Msg("caught signal")

		// Allow processes to finish with a ten-second window
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}
		app.Logger.Warn().Str("web-server", srv.Addr).Msg("completing background tasks")
		// Call wait so that the wait group can decrement to zero.
		wg.Wait()
		shutdownError <- nil
	}()

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-shutdownError
	if err != nil {
		app.Logger.Warn().Str("web-server", srv.Addr).Msg("stopped server")
		return err
	}
	return nil
}
