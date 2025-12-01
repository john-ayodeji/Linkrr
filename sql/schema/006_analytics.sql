-- +goose Up
CREATE TABLE analytics (
    id UUID PRIMARY KEY NOT NULL,
    short_code TEXT NOT NULL REFERENCES urls(short_code) ON DELETE CASCADE,
    clicked_at TIMESTAMP DEFAULT NOW(),
    ip inet NOT NULL,
    country TEXT NOT NULL,
    city TEXT NOT NULL,
    referrer TEXT NOT NULL,
    device TEXT NOT NULL,
    os TEXT NOT NULL,
    browser TEXT NOT NULL
);

-- +goose Down
DROP TABLE analytics;