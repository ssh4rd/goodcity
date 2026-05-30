-- +goose Up
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS is_role_verified BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS verification_doc_url TEXT;

-- +goose Down
ALTER TABLE users
    DROP COLUMN IF EXISTS is_role_verified,
    DROP COLUMN IF EXISTS verification_doc_url;
