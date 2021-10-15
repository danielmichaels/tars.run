package data

import (
	"database/sql"
	"errors"
)

// Model related errors
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Links     LinkModel
	Analytics AnalyticsModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Links:     LinkModel{DB: db},
		Analytics: AnalyticsModel{DB: db},
	}
}
