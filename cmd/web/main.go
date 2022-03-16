package main

import (
	"context"
	"database/sql"
	"github.com/danielmichaels/shortlink-go/internal/config"
	"github.com/danielmichaels/shortlink-go/internal/data"
	"github.com/danielmichaels/shortlink-go/internal/templates"
	"github.com/go-chi/httplog"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"html/template"
	"sync"
	"time"
)

type application struct {
	config   *config.Conf
	logger   zerolog.Logger
	template map[string]*template.Template
	models   data.Models
	wg       sync.WaitGroup
}

type templateData struct {
	Title   string
	AppName string
}

func main() {
	cfg := config.AppConfig()
	logger := httplog.NewLogger("web-server", httplog.Options{
		JSON:     cfg.Server.LogJson,
		Concise:  cfg.Server.LogConcise,
		LogLevel: cfg.Server.LogLevel,
	})
	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to open database. exiting")
	}
	logger.Info().Msg("database connection established")
	templateCache, err := templates.NewTemplateCache()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create a template cache")
	}
	app := &application{
		config:   cfg,
		logger:   logger,
		models:   data.NewModels(db),
		template: templateCache,
	}

	err = app.serve()
	if err != nil {
		app.logger.Error().Err(err).Msg("server failed to start")
	}
}

// openDB returns a sql connection pool
func openDB(cfg *config.Conf) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the
	// config struct
	db, err := sql.Open("sqlite3", cfg.Db.DbName)
	if err != nil {
		return nil, err
	}

	// Create a context with a 5-second timeout deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// PingContext establishes a new connection to the database, passing in the
	// ctx as a parameter. If the connection couldn't be established within
	// 5 seconds, an error will be raised.
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
