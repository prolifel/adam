-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS container_profiles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    collection_name TEXT NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(collection_name, key, value)
);

CREATE INDEX idx_collection_name ON container_profiles(collection_name);
CREATE INDEX idx_key ON container_profiles(key);
CREATE INDEX idx_collection_key ON container_profiles(collection_name, key);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_collection_key;
DROP INDEX IF EXISTS idx_key;
DROP INDEX IF EXISTS idx_collection_name;
DROP TABLE IF EXISTS container_profiles;
-- +goose StatementEnd
