-- +goose Up
CREATE TABLE IF NOT EXISTS fcm_tokens (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     TEXT         NOT NULL,
    token       TEXT         NOT NULL UNIQUE,
    platform    TEXT         NOT NULL CHECK (platform IN ('android', 'ios', 'web')),
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_fcm_tokens_user_id ON fcm_tokens(user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_fcm_tokens_user_id;
DROP TABLE IF EXISTS fcm_tokens;
