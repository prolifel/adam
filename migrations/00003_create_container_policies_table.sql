-- +goose Up
-- +goose StatementBegin
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

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_collection_name_policy;
DROP INDEX IF EXISTS idx_policy_id;
DROP TABLE IF EXISTS container_policies;
-- +goose StatementEnd
