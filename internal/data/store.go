package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/danielmichaels/shortlink-go/internal/config"
	_ "modernc.org/sqlite"
)

// OpenDB returns a sql connection pool
func OpenDB(cfg *config.Conf) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the
	// config struct
	db, err := sql.Open("sqlite", cfg.Db.DbName)
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
