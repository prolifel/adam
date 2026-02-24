-- +goose Up
-- +goose StatementBegin
ALTER TABLE container_profiles ADD COLUMN updated_at DATETIME;
ALTER TABLE container_profiles ADD COLUMN verdict TEXT;
ALTER TABLE container_profiles ADD COLUMN remarks TEXT;

CREATE INDEX idx_verdict ON container_profiles(verdict);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_verdict;
ALTER TABLE container_profiles DROP COLUMN remarks;
ALTER TABLE container_profiles DROP COLUMN verdict;
ALTER TABLE container_profiles DROP COLUMN updated_at;
-- +goose StatementEnd
