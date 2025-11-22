-- name: CreateUser :one
INSERT INTO users (
    username, email, password
)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 OR email = $2;