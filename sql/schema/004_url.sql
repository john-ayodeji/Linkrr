-- +goose Up
CREATE TABLE urls (
    id UUID PRIMARY KEY,
    short_code TEXT UNIQUE NOT NULL,
    url TEXT NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE urls;