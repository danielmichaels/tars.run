-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS analytics (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    ip TEXT NOT NULL DEFAULT 'no ip',
    user_agent TEXT NOT NULL,
    date_accessed DATE DEFAULT (datetime('now', 'utc')),
    links_id INTEGER NOT NULL,
    FOREIGN KEY (links_id) REFERENCES links (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS analytics;
-- +goose StatementEnd
