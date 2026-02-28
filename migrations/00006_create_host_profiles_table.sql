-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS host_profiles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    host_id TEXT NOT NULL,
    collection_name TEXT NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(host_id, collection_name, key, value)
);

CREATE INDEX idx_host_id ON host_profiles(host_id);
CREATE INDEX idx_host_collection_name ON host_profiles(collection_name);
CREATE INDEX idx_host_key ON host_profiles(key);
CREATE INDEX idx_host_collection_key ON host_profiles(collection_name, key);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_host_collection_key;
DROP INDEX IF EXISTS idx_host_key;
DROP INDEX IF EXISTS idx_host_collection_name;
DROP INDEX IF EXISTS idx_host_id;
DROP TABLE IF EXISTS host_profiles;
-- +goose StatementEnd
