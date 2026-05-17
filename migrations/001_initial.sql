-- +goose Up
CREATE TABLE users (
    id          BIGSERIAL PRIMARY KEY,
    email       TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role        TEXT NOT NULL DEFAULT 'user' CHECK (role IN ('user', 'moderator')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE practices (
    id          BIGSERIAL PRIMARY KEY,
    title       TEXT NOT NULL,
    description TEXT NOT NULL,
    city        TEXT NOT NULL,
    category    TEXT NOT NULL,
    status      TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    author_id   BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX practices_status_idx ON practices(status);
CREATE INDEX practices_city_idx ON practices(city);
CREATE INDEX practices_category_idx ON practices(category);

CREATE TABLE votes (
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    practice_id BIGINT NOT NULL REFERENCES practices(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, practice_id)
);

CREATE TABLE comments (
    id          BIGSERIAL PRIMARY KEY,
    practice_id BIGINT NOT NULL REFERENCES practices(id) ON DELETE CASCADE,
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    text        TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX comments_practice_idx ON comments(practice_id);

-- +goose Down
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS votes;
DROP TABLE IF EXISTS practices;
DROP TABLE IF EXISTS users;
