-- name: RefreshTokenCreate :one
INSERT INTO refresh_tokens (user_id, token, expires_at) 
VALUES (?, ?, ?) RETURNING id;

-- name: RefreshTokenDeleteByToken :one
DELETE FROM refresh_tokens WHERE token = ? RETURNING *;