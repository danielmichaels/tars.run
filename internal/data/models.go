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

func NewMockModels() Models {
	return Models{
		Links:     MockLinksModel{},
		Analytics: MockAnalyticsModel{},
	}
}
