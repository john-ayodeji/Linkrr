-- name: CreateAnalyticsData :exec
INSERT INTO analytics (
    id, short_code, alias, clicked_at, ip, country, referrer, device, os, browser
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);