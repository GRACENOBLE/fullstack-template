-- +goose Up
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS firebase_uid TEXT UNIQUE,
    ADD COLUMN IF NOT EXISTS email        TEXT,
    ADD COLUMN IF NOT EXISTS photo_url    TEXT,
    ADD COLUMN IF NOT EXISTS updated_at   TIMESTAMPTZ NOT NULL DEFAULT now();

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_firebase_uid ON users(firebase_uid);

-- +goose Down
DROP INDEX IF EXISTS idx_users_firebase_uid;
ALTER TABLE users
    DROP COLUMN IF EXISTS firebase_uid,
    DROP COLUMN IF EXISTS email,
    DROP COLUMN IF EXISTS photo_url,
    DROP COLUMN IF EXISTS updated_at;
