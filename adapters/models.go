package adapters

import (
	"fmt"
	"gorm.io/gorm"
	"shorty"
	"time"
)

// Link holds the domain information
// *note* we do not store the complete short link, instead we only supply the
// `hash`. This is so that we can prepend the DOMAIN with the hash during the
// translation - it prevents database corruption if the domain name is the change.
type Link struct {
	gorm.Model  `json:"-"`
	OriginalURL string       `json:"original_url"`
	Hash        string       `json:"hash"`
	Data        []DataPoints `json:"data,omitempty"`
}

// DataPoints belongs to Link,
type DataPoints struct {
	//gorm.Model             //`json:"-"`
	Id           uint      `json:"id" gorm:"primaryKey"`
	LinkID       string    `json:"link_id" gorm:"foreignKey:Hash"`
	Ip           string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	Location     string    `json:"location"`
	DateAccessed time.Time `json:"date_accessed"`
}

// CreateShortLink concatenates the current domain and the hash of the link
// it is created on the fly so that if the underlying domain changes the links
// are not broken in the future.
func (l *Link) CreateShortLink() string {
	return fmt.Sprintf("%s/%s", shorty.AppConfig().Server.AllowedOrigins, l.Hash)
}
