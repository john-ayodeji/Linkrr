-- +goose Up
CREATE TABLE password_reset (
    id UUID PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL,
    hashed_token TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN NOT NULL DEFAULT false
);

-- +goose Down
DROP TABLE password_reset;