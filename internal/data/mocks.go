package data

import (
	"log"
	"time"
)

type MockAnalyticsModel struct{}

func (m MockAnalyticsModel) GetAllForLink(hash string, filters Filters) ([]*Analytics, Metadata, error) {
	analytics := []*Analytics{
		{ID: 1, UserAgent: "test-agent", LinkID: 1, Ip: "1.1.1.1", DateAccessed: time.Date(2020, 01, 01, 00, 00, 00, 00, time.UTC)},
		{ID: 1, UserAgent: "test-agent", LinkID: 1, Ip: "1.1.1.1", DateAccessed: time.Date(2020, 01, 01, 00, 00, 00, 00, time.UTC)},
	}

	return analytics, Metadata{}, nil
}

func (m MockAnalyticsModel) Insert(analytics *Analytics) error {
	log.Print("inserting analytics")
	return nil
}

type MockLinksModel struct{}

var mockLink = &Link{
	ID:          1,
	CreatedAt:   time.Date(2020, 01, 01, 00, 00, 00, 00, time.UTC),
	OriginalURL: "test.com",
	Hash:        "notfake",
}

func (m MockLinksModel) Get(hash string) (*Link, error) {
	switch hash {
	case "notfake":
		return mockLink, nil
	default:
		return nil, ErrRecordNotFound
	}
}

func (m MockLinksModel) Insert(link *Link) error {
	return nil
}
