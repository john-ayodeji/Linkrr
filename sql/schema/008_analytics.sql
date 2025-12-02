-- +goose Up

-- 1. Index on short_code for fast filtering by URL
CREATE INDEX IF NOT EXISTS idx_analytics_short_code
    ON analytics(short_code);

-- 2. Composite index on short_code + alias for alias-level analytics
CREATE INDEX IF NOT EXISTS idx_analytics_short_code_alias
    ON analytics(short_code, alias);

-- 3. Index on clicked_at for time-based aggregation
CREATE INDEX IF NOT EXISTS idx_analytics_clicked_at
    ON analytics(clicked_at);

-- 4. Composite index for short_code + clicked_at (common for daily/weekly/monthly queries)
CREATE INDEX IF NOT EXISTS idx_analytics_short_code_clicked_at
    ON analytics(short_code, clicked_at);

-- 5. Index on IP for unique visitor calculations
CREATE INDEX IF NOT EXISTS idx_analytics_ip
    ON analytics(ip);

-- 6. Composite index on short_code + IP for per-URL unique visitors
CREATE INDEX IF NOT EXISTS idx_analytics_short_code_ip
    ON analytics(short_code, ip);

-- 7. Composite index for country + city grouping
CREATE INDEX IF NOT EXISTS idx_analytics_short_code_country_city
    ON analytics(short_code, country, city);

-- 8. Index for devices
CREATE INDEX IF NOT EXISTS idx_analytics_short_code_device
    ON analytics(short_code, device);

-- 9. Index for browsers
CREATE INDEX IF NOT EXISTS idx_analytics_short_code_browser
    ON analytics(short_code, browser);

-- 10. Index for referrer
CREATE INDEX IF NOT EXISTS idx_analytics_short_code_referrer
    ON analytics(short_code, referrer);

-- 11. Index for urls.user_id to speed up global analytics per user
CREATE INDEX IF NOT EXISTS idx_urls_user_id
    ON urls(user_id);

-- 12. Optional covering index for the most common analytics queries
--    (reduces table scans for analytics dashboards)
CREATE INDEX IF NOT EXISTS idx_analytics_covering
    ON analytics(short_code, alias, clicked_at, country, city, device, browser, referrer, ip);
