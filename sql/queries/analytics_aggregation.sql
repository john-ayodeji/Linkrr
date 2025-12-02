-- URL-level aggregation queries

-- name: UpsertURLDaily :exec
INSERT INTO analytics_url_daily (short_code, date, total_clicks, unique_visitors)
VALUES ($1, $2, $3, $4)
ON CONFLICT (short_code, date)
DO UPDATE SET
    total_clicks = analytics_url_daily.total_clicks + EXCLUDED.total_clicks,
    unique_visitors = EXCLUDED.unique_visitors;

-- name: UpsertURLGeo :exec
INSERT INTO analytics_url_geo (short_code, country, city, total_clicks)
VALUES ($1, $2, $3, $4)
ON CONFLICT (short_code, country, city)
DO UPDATE SET
    total_clicks = analytics_url_geo.total_clicks + EXCLUDED.total_clicks;

-- name: UpsertURLDevice :exec
INSERT INTO analytics_url_device (short_code, device, total_clicks)
VALUES ($1, $2, $3)
ON CONFLICT (short_code, device)
DO UPDATE SET
    total_clicks = analytics_url_device.total_clicks + EXCLUDED.total_clicks;

-- name: UpsertURLBrowser :exec
INSERT INTO analytics_url_browser (short_code, browser, total_clicks)
VALUES ($1, $2, $3)
ON CONFLICT (short_code, browser)
DO UPDATE SET
    total_clicks = analytics_url_browser.total_clicks + EXCLUDED.total_clicks;

-- name: UpsertURLReferrer :exec
INSERT INTO analytics_url_referrer (short_code, referrer, total_clicks)
VALUES ($1, $2, $3)
ON CONFLICT (short_code, referrer)
DO UPDATE SET
    total_clicks = analytics_url_referrer.total_clicks + EXCLUDED.total_clicks;

-- Alias-level aggregation queries

-- name: UpsertAliasDaily :exec
INSERT INTO analytics_alias_daily (short_code, alias, date, total_clicks, unique_visitors)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (short_code, alias, date)
DO UPDATE SET
    total_clicks = analytics_alias_daily.total_clicks + EXCLUDED.total_clicks,
    unique_visitors = EXCLUDED.unique_visitors;

-- name: UpsertAliasGeo :exec
INSERT INTO analytics_alias_geo (short_code, alias, country, city, total_clicks)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (short_code, alias, country, city)
DO UPDATE SET
    total_clicks = analytics_alias_geo.total_clicks + EXCLUDED.total_clicks;

-- name: UpsertAliasDevice :exec
INSERT INTO analytics_alias_device (short_code, alias, device, total_clicks)
VALUES ($1, $2, $3, $4)
ON CONFLICT (short_code, alias, device)
DO UPDATE SET
    total_clicks = analytics_alias_device.total_clicks + EXCLUDED.total_clicks;

-- name: UpsertAliasBrowser :exec
INSERT INTO analytics_alias_browser (short_code, alias, browser, total_clicks)
VALUES ($1, $2, $3, $4)
ON CONFLICT (short_code, alias, browser)
DO UPDATE SET
    total_clicks = analytics_alias_browser.total_clicks + EXCLUDED.total_clicks;

-- name: UpsertAliasReferrer :exec
INSERT INTO analytics_alias_referrer (short_code, alias, referrer, total_clicks)
VALUES ($1, $2, $3, $4)
ON CONFLICT (short_code, alias, referrer)
DO UPDATE SET
    total_clicks = analytics_alias_referrer.total_clicks + EXCLUDED.total_clicks;

-- User-level aggregation queries

-- name: UpsertUserDaily :exec
INSERT INTO analytics_user_daily (user_id, date, total_clicks, unique_visitors)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id, date)
DO UPDATE SET
    total_clicks = analytics_user_daily.total_clicks + EXCLUDED.total_clicks,
    unique_visitors = EXCLUDED.unique_visitors;

-- name: UpsertUserGeo :exec
INSERT INTO analytics_user_geo (user_id, country, total_clicks)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, country)
DO UPDATE SET
    total_clicks = analytics_user_geo.total_clicks + EXCLUDED.total_clicks;

-- name: UpsertUserBrowser :exec
INSERT INTO analytics_user_browser (user_id, browser, total_clicks)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, browser)
DO UPDATE SET
    total_clicks = analytics_user_browser.total_clicks + EXCLUDED.total_clicks;

-- name: UpsertUserTopLinks :exec
INSERT INTO analytics_user_top_links (user_id, short_code, total_clicks)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, short_code)
DO UPDATE SET
    total_clicks = analytics_user_top_links.total_clicks + EXCLUDED.total_clicks;

-- Helper queries to get aggregated data for a single click event

-- name: GetClickEventData :one
SELECT
    a.short_code,
    a.alias,
    a.clicked_at,
    a.ip,
    a.country,
    a.city,
    a.referrer,
    a.device,
    a.browser,
    u.user_id
FROM analytics a
JOIN urls u ON u.short_code = a.short_code
WHERE a.short_code = $1
  AND a.clicked_at = $2
  AND a.ip = $3
LIMIT 1;

-- name: GetURLTotalClicks :one
SELECT COALESCE(SUM(total_clicks), 0) as total_clicks
FROM analytics_url_daily
WHERE short_code = $1;

-- name: GetURLUniqueVisitors :one
SELECT COALESCE(SUM(unique_visitors), 0) as unique_visitors
FROM analytics_url_daily
WHERE short_code = $1;

-- name: GetURLDailyClicks :many
SELECT date, total_clicks, unique_visitors
FROM analytics_url_daily
WHERE short_code = $1
ORDER BY date DESC
LIMIT 30;

-- name: GetURLGeoClicks :many
SELECT country, city, total_clicks
FROM analytics_url_geo
WHERE short_code = $1
ORDER BY total_clicks DESC
LIMIT 10;

-- name: GetURLDeviceClicks :many
SELECT device, total_clicks
FROM analytics_url_device
WHERE short_code = $1
ORDER BY total_clicks DESC;

-- name: GetURLBrowserClicks :many
SELECT browser, total_clicks
FROM analytics_url_browser
WHERE short_code = $1
ORDER BY total_clicks DESC
LIMIT 10;

-- name: GetURLReferrerClicks :many
SELECT referrer, total_clicks
FROM analytics_url_referrer
WHERE short_code = $1
ORDER BY total_clicks DESC
LIMIT 10;

-- Alias-level aggregate reads

-- name: GetAliasTotalClicks :one
SELECT COALESCE(SUM(total_clicks), 0) as total_clicks
FROM analytics_alias_daily
WHERE short_code = $1 AND alias = $2;

-- name: GetAliasUniqueVisitors :one
SELECT COALESCE(SUM(unique_visitors), 0) as unique_visitors
FROM analytics_alias_daily
WHERE short_code = $1 AND alias = $2;

-- name: GetAliasDailyClicks :many
SELECT date, total_clicks, unique_visitors
FROM analytics_alias_daily
WHERE short_code = $1 AND alias = $2
ORDER BY date DESC
LIMIT 30;

-- name: GetAliasGeoClicks :many
SELECT country, city, total_clicks
FROM analytics_alias_geo
WHERE short_code = $1 AND alias = $2
ORDER BY total_clicks DESC
LIMIT 10;

-- name: GetAliasDeviceClicks :many
SELECT device, total_clicks
FROM analytics_alias_device
WHERE short_code = $1 AND alias = $2
ORDER BY total_clicks DESC;

-- name: GetAliasBrowserClicks :many
SELECT browser, total_clicks
FROM analytics_alias_browser
WHERE short_code = $1 AND alias = $2
ORDER BY total_clicks DESC
LIMIT 10;

-- name: GetAliasReferrerClicks :many
SELECT referrer, total_clicks
FROM analytics_alias_referrer
WHERE short_code = $1 AND alias = $2
ORDER BY total_clicks DESC
LIMIT 10;

-- User-level aggregate reads

-- name: GetUserTotalClicks :one
SELECT COALESCE(SUM(total_clicks), 0) as total_clicks
FROM analytics_user_daily
WHERE user_id = $1;

-- name: GetUserUniqueVisitors :one
SELECT COALESCE(SUM(unique_visitors), 0) as unique_visitors
FROM analytics_user_daily
WHERE user_id = $1;

-- name: GetUserDailyClicks :many
SELECT date, total_clicks, unique_visitors
FROM analytics_user_daily
WHERE user_id = $1
ORDER BY date DESC
LIMIT 30;

-- name: GetUserGeoClicks :many
SELECT country, total_clicks
FROM analytics_user_geo
WHERE user_id = $1
ORDER BY total_clicks DESC
LIMIT 10;

-- name: GetUserBrowserClicks :many
SELECT browser, total_clicks
FROM analytics_user_browser
WHERE user_id = $1
ORDER BY total_clicks DESC
LIMIT 10;

-- name: GetUserTopLinks :many
SELECT short_code, total_clicks
FROM analytics_user_top_links
WHERE user_id = $1
ORDER BY total_clicks DESC
LIMIT 10;
