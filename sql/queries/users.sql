-- name: CreateUser :one
INSERT INTO users (
    username, email, password
)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 OR email = $2;

-- name: UpdatePassword :one
UPDATE users
SET password = $1
WHERE id = $2
RETURNING *;