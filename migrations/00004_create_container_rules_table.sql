-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS container_policies;

CREATE TABLE IF NOT EXISTS container_rules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    collection_name TEXT NOT NULL,
    rule TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(collection_name)
);

CREATE INDEX idx_collection_name_rules ON container_rules(collection_name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_collection_name_rules;
DROP TABLE IF EXISTS container_rules;

CREATE TABLE IF NOT EXISTS container_policies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    policy_id TEXT NOT NULL,
    collection_name TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(policy_id, collection_name)
);

CREATE INDEX idx_policy_id ON container_policies(policy_id);
CREATE INDEX idx_collection_name_policy ON container_policies(collection_name);
-- +goose StatementEnd
