-- +goose Up
CREATE TABLE aliases (
    id UUID PRIMARY KEY,
    alias TEXT UNIQUE NOT NULL,
    url_code TEXT NOT NULL REFERENCES urls(short_code) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE aliases;