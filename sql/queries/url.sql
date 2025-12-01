-- name: CreateURL :one
INSERT INTO urls (
    id, short_code, url, user_id
)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetURL :one
SELECT * FROM urls
WHERE short_code = $1;

-- name: GetOriginalUrl :one
SELECT urls.url AS link FROM urls
LEFT JOIN aliases
    ON aliases.url_code = urls.short_code
WHERE aliases.alias = $1
   OR urls.short_code = $2;

-- name: GetShortCodeAndAlias :one
SELECT urls.short_code AS code, aliases.alias AS alias
FROM urls
         LEFT JOIN aliases ON aliases.url_code = urls.short_code
WHERE $1 IN (aliases.alias, urls.short_code)
LIMIT 1;
