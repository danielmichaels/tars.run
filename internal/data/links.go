package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/danielmichaels/shortlink-go/internal/validator"
)

// Link holds the domain information
// *note* we do not store the complete short link, instead we only supply the
// `hash`. This is so that we can prepend the DOMAIN with the hash during the
// translation - it prevents database corruption if the domain name is the change.
type Link struct {
	CreatedAt   time.Time `json:"created_at"` // todo: omit?
	OriginalURL string    `json:"original_url"`
	Hash        string    `json:"hash"`
	ID          int64     `json:"id"` // todo: omit
}

// CreateShortLink concatenates the current domain and the hash of the link
// it is created on the fly so that if the underlying domain changes the links
// are not broken in the future.
func (l *Link) CreateShortLink() string {
	return fmt.Sprintf("%s/%s", os.Getenv("DOMAIN"), l.Hash)
}

type LinkModel struct {
	DB *sql.DB
}

// Insert a new link into the database.
func (m LinkModel) Insert(link *Link) error {
	query := `
		INSERT INTO links (original_url, hash)
		VALUES ($1, $2)
		RETURNING id, created_at`
	// RETURNING supported by sqlite3.35+

	args := []interface{}{link.OriginalURL, link.Hash}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&link.ID, &link.CreatedAt)
}

// Get retrieves a specific link
func (m LinkModel) Get(hash string) (*Link, error) {
	// todo validate on hash string (need to make a hash type)
	query := `
		SELECT id, created_at, original_url, hash
		FROM links
		WHERE hash = $1`
	var link Link

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, hash).Scan(
		&link.ID,
		&link.CreatedAt,
		&link.OriginalURL,
		&link.Hash,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &link, nil
}

// ValidateURL checks that the URL starts with a valid scheme
func ValidateURL(v *validator.Validator, url string) {
	v.Check(validator.IsURL(url), "link", "url must start with a valid scheme such as https://")
}
