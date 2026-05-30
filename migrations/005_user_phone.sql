-- +goose Up
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone TEXT;

-- +goose Down
ALTER TABLE users DROP COLUMN IF EXISTS phone;
