-- name: CreateAnalyticsData :one
INSERT INTO analytics (
    id, short_code, alias, clicked_at, ip, country, city, referrer, device, os, browser
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: TotalClicksPerURL :one
SELECT COUNT(*) AS total_clicks
FROM analytics
         JOIN urls ON urls.short_code = analytics.short_code
WHERE analytics.short_code = $1
  AND urls.user_id = $2;

-- name: UniqueVisitorsPerURL :one
SELECT COUNT(DISTINCT ip) AS unique_visitors
FROM analytics
         JOIN urls ON urls.short_code = analytics.short_code
WHERE analytics.short_code = $1
  AND urls.user_id = $2;

-- name: DailyClicksPerURL :many
SELECT DATE(clicked_at) AS day, COUNT(*) AS clicks
FROM analytics
         JOIN urls ON urls.short_code = analytics.short_code
WHERE analytics.short_code = $1
  AND urls.user_id = $2
GROUP BY day
ORDER BY day;

-- name: WeeklyClicksPerURL :many
SELECT DATE_TRUNC('week', clicked_at) AS week, COUNT(*) AS clicks
FROM analytics
         JOIN urls ON urls.short_code = analytics.short_code
WHERE analytics.short_code = $1
  AND urls.user_id = $2
GROUP BY week
ORDER BY week;

-- name: MonthlyClicksPerURL :many
SELECT DATE_TRUNC('month', clicked_at) AS month, COUNT(*) AS clicks
FROM analytics
         JOIN urls ON urls.short_code = analytics.short_code
WHERE analytics.short_code = $1
  AND urls.user_id = $2
GROUP BY month
ORDER BY month;

-- name: CountryClicksPerURL :many
SELECT country, COUNT(*) AS clicks
FROM analytics
         JOIN urls ON urls.short_code = analytics.short_code
WHERE analytics.short_code = $1
  AND urls.user_id = $2
GROUP BY country
ORDER BY clicks DESC;

-- name: CityClicksPerURL :many
SELECT city, COUNT(*) AS clicks
FROM analytics
         JOIN urls ON urls.short_code = analytics.short_code
WHERE analytics.short_code = $1
  AND urls.user_id = $2
GROUP BY city
ORDER BY clicks DESC;

-- name: DeviceClicksPerURL :many
SELECT device, COUNT(*) AS clicks
FROM analytics
         JOIN urls ON urls.short_code = analytics.short_code
WHERE analytics.short_code = $1
  AND urls.user_id = $2
GROUP BY device
ORDER BY clicks DESC;

-- name: BrowserClicksPerURL :many
SELECT browser, COUNT(*) AS clicks
FROM analytics
         JOIN urls ON urls.short_code = analytics.short_code
WHERE analytics.short_code = $1
  AND urls.user_id = $2
GROUP BY browser
ORDER BY clicks DESC;

-- name: RefererPerURL :many
SELECT COALESCE(referrer, 'direct') AS referrer, COUNT(*) AS clicks
FROM analytics
         JOIN urls ON urls.short_code = analytics.short_code
WHERE analytics.short_code = $1
  AND urls.user_id = $2
GROUP BY referrer
ORDER BY clicks DESC;

-- name: ClicksByRefererPerURL :many
SELECT COALESCE(referrer, 'direct') AS referrer, COUNT(*) AS total
FROM analytics
         JOIN urls ON urls.short_code = analytics.short_code
WHERE analytics.short_code = $1
  AND urls.user_id = $2
GROUP BY referrer
ORDER BY total DESC;

-- name: TotalClicksPerAlias :one
SELECT COUNT(*) AS total
FROM analytics
         JOIN urls ON urls.short_code = analytics.short_code
WHERE analytics.short_code = $1
  AND urls.user_id = $2
  AND analytics.alias = $3;

-- name: UniqueVisitorsPerAlias :one
SELECT COUNT(DISTINCT ip) AS unique_visitors
FROM analytics
         JOIN urls ON urls.short_code = analytics.short_code
WHERE analytics.short_code = $1
  AND urls.user_id = $2
  AND analytics.alias = $3;

-- name: DailyClicksPerAlias :many
SELECT DATE(clicked_at) AS day, COUNT(*) AS total
FROM analytics
         JOIN urls ON urls.short_code = analytics.short_code
WHERE analytics.short_code = $1
  AND urls.user_id = $2
  AND analytics.alias = $3
GROUP BY day
ORDER BY day;

-- name: CountryClicksPerAlias :many
SELECT country, COUNT(*) AS clicks
FROM analytics
         JOIN urls ON urls.short_code = analytics.short_code
WHERE analytics.short_code = $1
  AND urls.user_id = $2
  AND analytics.alias = $3
GROUP BY country
ORDER BY clicks DESC;

-- name: CityClicksPerAlias :many
SELECT city, COUNT(*) AS clicks
FROM analytics
         JOIN urls ON urls.short_code = analytics.short_code
WHERE analytics.short_code = $1
  AND urls.user_id = $2
  AND analytics.alias = $3
GROUP BY city
ORDER BY clicks DESC;

-- name: DeviceClicksPerAlias :many
SELECT device, COUNT(*) AS clicks
FROM analytics
         JOIN urls ON urls.short_code = analytics.short_code
WHERE analytics.short_code = $1
  AND urls.user_id = $2
  AND analytics.alias = $3
GROUP BY device
ORDER BY clicks DESC;

-- name: BrowserClicksPerAlias :many
SELECT browser, COUNT(*) AS clicks
FROM analytics
         JOIN urls ON urls.short_code = analytics.short_code
WHERE analytics.short_code = $1
  AND urls.user_id = $2
  AND analytics.alias = $3
GROUP BY browser
ORDER BY clicks DESC;

-- name: RefererPerAlias :many
SELECT COALESCE(referrer, 'direct') AS referrer, COUNT(*) AS clicks
FROM analytics
         JOIN urls ON urls.short_code = analytics.short_code
WHERE analytics.short_code = $1
  AND urls.user_id = $2
  AND analytics.alias = $3
GROUP BY referrer
ORDER BY clicks DESC;

-- name: TotalClicksGlobal :one
SELECT COUNT(*) AS total
FROM analytics
WHERE short_code IN (
    SELECT short_code FROM urls WHERE user_id = $1
);

-- name: UniqueVisitorsGlobal :one
SELECT COUNT(DISTINCT ip) AS unique_visitors
FROM analytics
WHERE short_code IN (
    SELECT short_code FROM urls WHERE user_id = $1
);

-- name: DailyClicksGlobal :many
SELECT DATE(clicked_at) AS day, COUNT(*) AS total
FROM analytics
WHERE short_code IN (
    SELECT short_code FROM urls WHERE user_id = $1
)
GROUP BY day
ORDER BY day;

-- name: CountryClicksGlobal :many
SELECT country, COUNT(*) AS clicks
FROM analytics
WHERE short_code IN (
    SELECT short_code FROM urls WHERE user_id = $1
)
GROUP BY country
ORDER BY clicks DESC;

-- name: BrowserClicksGlobal :many
SELECT browser, COUNT(*) AS clicks
FROM analytics
WHERE short_code IN (
    SELECT short_code FROM urls WHERE user_id = $1
)
GROUP BY browser
ORDER BY clicks DESC;

-- name: TopPerformingLinksPerUser :many
SELECT short_code, COUNT(*) AS total
FROM analytics
WHERE short_code IN (
    SELECT short_code FROM urls WHERE user_id = $1
)
GROUP BY short_code
ORDER BY total DESC;