-- name: CreateToken :exec
INSERT INTO password_reset (
    id, user_id, hashed_token, expires_at
)
VALUES ($1, $2, $3, $4);

-- name: GetToken :one
SELECT user_id, expires_at, used FROM password_reset
WHERE hashed_token = $1;

-- name: SetUsed :exec
UPDATE password_reset
SET used = true;