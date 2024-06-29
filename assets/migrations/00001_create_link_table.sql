-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS links (
    id INTEGER NOT NULL primary key autoincrement,
    created_at DATE DEFAULT (datetime('now', 'utc')),
    original_url TEXT NOT NULL,
    hash TEXT NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS links;
-- +goose StatementEnd
