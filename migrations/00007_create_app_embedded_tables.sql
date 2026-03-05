-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS app_embedded_profiles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    profile_id TEXT NOT NULL,
    app_id TEXT,
    collection_name TEXT NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(profile_id, collection_name, key, value)
);

CREATE TABLE IF NOT EXISTS app_embedded_policies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    policy_id TEXT NOT NULL,
    collection_name TEXT NOT NULL,
    rule TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(policy_id, collection_name)
);

CREATE INDEX idx_app_embedded_profile_id ON app_embedded_profiles(profile_id);
CREATE INDEX idx_app_embedded_collection ON app_embedded_profiles(collection_name);
CREATE INDEX idx_app_embedded_key ON app_embedded_profiles(key);
CREATE INDEX idx_app_embedded_policy_id ON app_embedded_policies(policy_id);
CREATE INDEX idx_app_embedded_policy_collection ON app_embedded_policies(collection_name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_app_embedded_policy_collection;
DROP INDEX IF EXISTS idx_app_embedded_policy_id;
DROP INDEX IF EXISTS idx_app_embedded_key;
DROP INDEX IF EXISTS idx_app_embedded_collection;
DROP INDEX IF EXISTS idx_app_embedded_profile_id;
DROP TABLE IF EXISTS app_embedded_policies;
DROP TABLE IF EXISTS app_embedded_profiles;
-- +goose StatementEnd
