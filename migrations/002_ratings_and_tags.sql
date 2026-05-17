-- +goose Up

-- Extend users table
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS income_bracket TEXT     DEFAULT 'middle' CHECK (income_bracket IN ('low', 'middle', 'high')),
    ADD COLUMN IF NOT EXISTS social_role    TEXT     DEFAULT 'resident' CHECK (social_role IN ('resident', 'activist', 'specialist', 'official')),
    ADD COLUMN IF NOT EXISTS reputation_score DECIMAL(5,3) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS profile_weight   DECIMAL(5,3) NOT NULL DEFAULT 1;

-- Extend practices table
ALTER TABLE practices
    ADD COLUMN IF NOT EXISTS address        TEXT,
    ADD COLUMN IF NOT EXISTS latitude       DECIMAL(9,6),
    ADD COLUMN IF NOT EXISTS longitude      DECIMAL(9,6),
    ADD COLUMN IF NOT EXISTS budget_rub     DECIMAL(12,2),
    ADD COLUMN IF NOT EXISTS implemented_at DATE,
    ADD COLUMN IF NOT EXISTS rating         DECIMAL(4,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS rating_count   INT          NOT NULL DEFAULT 0;

-- Extend status constraint to include best_practice
ALTER TABLE practices DROP CONSTRAINT IF EXISTS practices_status_check;
ALTER TABLE practices
    ADD CONSTRAINT practices_status_check
    CHECK (status IN ('pending', 'approved', 'rejected', 'best_practice'));

-- Ratings table (replaces binary votes with multi-criteria scoring)
CREATE TABLE IF NOT EXISTS ratings (
    id               BIGSERIAL PRIMARY KEY,
    user_id          BIGINT       NOT NULL REFERENCES users(id)      ON DELETE CASCADE,
    practice_id      BIGINT       NOT NULL REFERENCES practices(id)  ON DELETE CASCADE,
    score_utility    SMALLINT     NOT NULL CHECK (score_utility    BETWEEN 1 AND 5),
    score_aesthetics SMALLINT     NOT NULL CHECK (score_aesthetics BETWEEN 1 AND 5),
    score_ecology    SMALLINT     NOT NULL CHECK (score_ecology    BETWEEN 1 AND 5),
    comment          TEXT,
    user_lat         DECIMAL(9,6),
    user_lon         DECIMAL(9,6),
    is_suspicious    BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, practice_id)
);

CREATE INDEX IF NOT EXISTS ratings_practice_idx ON ratings(practice_id);
CREATE INDEX IF NOT EXISTS ratings_user_idx     ON ratings(user_id);

-- Migrate existing votes → ratings with neutral scores (3,3,3)
INSERT INTO ratings (user_id, practice_id, score_utility, score_aesthetics, score_ecology)
SELECT user_id, practice_id, 3, 3, 3 FROM votes
ON CONFLICT (user_id, practice_id) DO NOTHING;

-- Tags
CREATE TABLE IF NOT EXISTS tags (
    id   SERIAL      PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    slug VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS practice_tags (
    practice_id BIGINT NOT NULL REFERENCES practices(id) ON DELETE CASCADE,
    tag_id      INT    NOT NULL REFERENCES tags(id)      ON DELETE CASCADE,
    PRIMARY KEY (practice_id, tag_id)
);

-- +goose Down
DROP TABLE IF EXISTS practice_tags;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS ratings;

ALTER TABLE practices DROP CONSTRAINT IF EXISTS practices_status_check;
ALTER TABLE practices ADD CONSTRAINT practices_status_check
    CHECK (status IN ('pending', 'approved', 'rejected'));
ALTER TABLE practices
    DROP COLUMN IF EXISTS address,
    DROP COLUMN IF EXISTS latitude,
    DROP COLUMN IF EXISTS longitude,
    DROP COLUMN IF EXISTS budget_rub,
    DROP COLUMN IF EXISTS implemented_at,
    DROP COLUMN IF EXISTS rating,
    DROP COLUMN IF EXISTS rating_count;

ALTER TABLE users
    DROP COLUMN IF EXISTS income_bracket,
    DROP COLUMN IF EXISTS social_role,
    DROP COLUMN IF EXISTS reputation_score,
    DROP COLUMN IF EXISTS profile_weight;
