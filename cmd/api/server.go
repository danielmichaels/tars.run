package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) serve() error {

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.Server.Port),
		Handler:      app.routes(),
		IdleTimeout:  app.config.Server.TimeoutIdle,
		ReadTimeout:  app.config.Server.TimeoutRead,
		WriteTimeout: app.config.Server.TimeoutWrite,
	}

	app.logger.PrintInfo("starting server", map[string]string{
		"addr":  srv.Addr,
		"debug": strconv.FormatBool(app.config.Debug),
	})

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
