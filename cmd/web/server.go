package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.Server.Port),
		Handler:      app.routes(),
		IdleTimeout:  app.config.Server.TimeoutIdle,
		ReadTimeout:  app.config.Server.TimeoutRead,
		WriteTimeout: app.config.Server.TimeoutWrite,
	}

	shutdownError := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.logger.Warn().Str("signal", s.String()).Msg("caught signal")

		// Allow processes to finish with a ten-second window
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}
		app.logger.Warn().Str("tasks", srv.Addr).Msg("completing background tasks")
		// Call wait so that the wait group can decrement to zero.
		app.wg.Wait()
		shutdownError <- nil
	}()
	app.logger.Info().Str("server", srv.Addr).Msg("starting server")

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-shutdownError
	if err != nil {
		return err
	}
	app.logger.Warn().Str("server", srv.Addr).Msg("stopped server")
	return nil
}
