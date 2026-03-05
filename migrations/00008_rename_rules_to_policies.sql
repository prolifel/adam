-- +goose Up
-- +goose StatementBegin
-- Drop container_rules and create container_policies
DROP TABLE IF EXISTS container_rules;

CREATE TABLE IF NOT EXISTS container_policies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    collection_name TEXT NOT NULL,
    rule TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(collection_name)
);

CREATE INDEX idx_collection_name_container_policies ON container_policies(collection_name);

-- Drop host_rules and create host_policies
DROP TABLE IF EXISTS host_rules;

CREATE TABLE IF NOT EXISTS host_policies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    collection_name TEXT NOT NULL,
    rule TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(collection_name)
);

CREATE INDEX idx_collection_name_host_policies ON host_policies(collection_name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop new tables and recreate old ones
DROP INDEX IF EXISTS idx_collection_name_host_policies;
DROP TABLE IF EXISTS host_policies;

DROP INDEX IF EXISTS idx_collection_name_container_policies;
DROP TABLE IF EXISTS container_policies;

CREATE TABLE IF NOT EXISTS host_rules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    collection_name TEXT NOT NULL,
    rule TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(collection_name)
);

CREATE INDEX idx_collection_name_host_rules ON host_rules(collection_name);

CREATE TABLE IF NOT EXISTS container_rules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    collection_name TEXT NOT NULL,
    rule TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(collection_name)
);

CREATE INDEX idx_collection_name_rules ON container_rules(collection_name);
-- +goose StatementEnd
