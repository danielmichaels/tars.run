package main

import (
	"context"
	"database/sql"
	"github.com/danielmichaels/shortlink-go/cmd/api/config"
	"github.com/danielmichaels/shortlink-go/internal/data"
	"github.com/danielmichaels/shortlink-go/internal/logger"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"sync"
	"time"
)

type application struct {
	config *config.Conf
	logger *logger.Logger
	models data.Models
	wg     sync.WaitGroup
}

func main() {
	cfg := config.AppConfig()
	log := logger.New(os.Stdout, logger.LevelInfo)
	db, err := openDB(cfg)
	if err != nil {
		log.PrintFatal(err, nil)
	}

	log.PrintInfo("database connection established", nil)

	app := &application{
		config: cfg,
		logger: log,
		models: data.NewModels(db),
	}

	err = app.serve()
	if err != nil {
		log.PrintFatal(err, nil)
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
