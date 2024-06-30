package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type LinkStore interface {
	Get(hash string) (*Link, error)
	Insert(link *Link) error
}
type AnalyticsStore interface {
	Insert(analytics *Analytics) error
	GetAllForLink(hash string, filters Filters) ([]*Analytics, Metadata, error)
}

type Models struct {
	Links     LinkStore
	Analytics AnalyticsStore
}

func NewModels(db *sql.DB) Models {
	return Models{
		Links:     LinkModel{DB: db},
		Analytics: AnalyticsModel{DB: db},
	}
}
