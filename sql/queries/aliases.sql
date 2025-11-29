-- name: CreateAlias :one
INSERT INTO aliases (
    id, alias, url_code
)
VALUES ($1, $2, $3)
RETURNING *;