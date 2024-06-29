package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Analytics belongs to Link,
type Analytics struct {
	//Location     string    `json:"location"` todo: add topo location data from IP
	DateAccessed time.Time `json:"date_accessed"`
	Ip           string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	ID           uint      `json:"id"`
	LinkID       uint64    `json:"-"`
}

type AnalyticsModel struct {
	DB *sql.DB
}

// Insert will create new entry for the Analytics table. This is typically
// called from the New method
func (m AnalyticsModel) Insert(analytics *Analytics) error {
	query := `
	INSERT INTO analytics (links_id, ip, user_agent)
	VALUES ($1, $2, $3)
	RETURNING date_accessed`
	// RETURNING supported by sqlite3.35+
	args := []interface{}{analytics.LinkID, analytics.Ip, analytics.UserAgent}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&analytics.DateAccessed,
	)
}

func (m AnalyticsModel) GetAllForLink(
	hash string,
	filters Filters,
) ([]*Analytics, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) over(), analytics.id, analytics.date_accessed, analytics.ip, analytics.user_agent
		FROM analytics
		INNER JOIN links
		ON analytics.links_id = links.id
		WHERE links.hash = $1
		ORDER BY %s %s, links.id ASC
		LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	var analytics []*Analytics

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{hash, filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	for rows.Next() {
		var analytic Analytics
		err := rows.Scan(
			&totalRecords,
			&analytic.ID,
			&analytic.DateAccessed,
			&analytic.Ip,
			&analytic.UserAgent,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		analytics = append(analytics, &analytic)
	}
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return analytics, metadata, nil
}
