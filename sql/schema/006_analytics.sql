-- +goose Up
CREATE TABLE analytics (
    id UUID PRIMARY KEY NOT NULL,
    short_code TEXT NOT NULL REFERENCES urls(short_code) ON DELETE CASCADE,
    alias TEXT REFERENCES aliases(alias),
    clicked_at TIMESTAMP NOT NULL,
    ip TEXT NOT NULL,
    country TEXT NOT NULL,
    referrer TEXT,
    device TEXT NOT NULL,
    os TEXT NOT NULL,
    browser TEXT NOT NULL
);

-- +goose Down
DROP TABLE analytics;