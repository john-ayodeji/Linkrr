-- name: CreateRefreshToken :exec
INSERT INTO refresh_token (
    token, user_id, expires_at
)
VALUES ($1, $2, $3);

-- name: VerifyRefreshToken :one
SELECT * FROM refresh_token
WHERE token = $1 AND revoked_at IS NULL;

-- name: RevokeRefreshToken :exec
UPDATE refresh_token
SET revoked_at = NOW()
WHERE token = $1;