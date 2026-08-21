-- name: RefreshTokenCreate :one
INSERT INTO refresh_tokens (user_id, token, expires_at) 
VALUES (?, ?, ?) RETURNING id;

-- name: RefreshTokenDeleteByID :one
DELETE FROM refresh_tokens WHERE id = ? RETURNING *;