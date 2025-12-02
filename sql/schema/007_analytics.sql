-- +goose Up
ALTER TABLE analytics
ADD COLUMN city TEXT NOT NULL DEFAULT 'lagos';