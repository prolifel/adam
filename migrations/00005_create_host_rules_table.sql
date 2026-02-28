-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS host_rules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    collection_name TEXT NOT NULL,
    rule TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(collection_name)
);

CREATE INDEX idx_collection_name_host_rules ON host_rules(collection_name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_collection_name_host_rules;
DROP TABLE IF EXISTS host_rules;
-- +goose StatementEnd
