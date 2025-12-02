-- name: CreateAlias :one
INSERT INTO aliases (
    id, alias, url_code
)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUrlOwnerByAlias :one
SELECT u.user_id
FROM aliases a
         JOIN urls u ON u.short_code = a.url_code
WHERE a.alias = $1
LIMIT 1;
