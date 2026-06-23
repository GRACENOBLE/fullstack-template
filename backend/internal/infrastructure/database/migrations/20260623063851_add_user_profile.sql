-- +goose Up
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS firebase_uid TEXT,
    ADD COLUMN IF NOT EXISTS email        TEXT,
    ADD COLUMN IF NOT EXISTS photo_url    TEXT,
    ADD COLUMN IF NOT EXISTS updated_at   TIMESTAMPTZ NOT NULL DEFAULT now();

-- Backfill any pre-existing rows so the NOT NULL constraint can be applied.
UPDATE users SET firebase_uid = 'legacy-' || id::text WHERE firebase_uid IS NULL;

ALTER TABLE users ALTER COLUMN firebase_uid SET NOT NULL;
ALTER TABLE users ADD CONSTRAINT users_firebase_uid_key UNIQUE (firebase_uid);

-- +goose Down
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_firebase_uid_key;
ALTER TABLE users
    DROP COLUMN IF EXISTS firebase_uid,
    DROP COLUMN IF EXISTS email,
    DROP COLUMN IF EXISTS photo_url,
    DROP COLUMN IF EXISTS updated_at;
