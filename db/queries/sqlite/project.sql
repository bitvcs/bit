-- name: CreateProject :one
INSERT INTO projects (name, description) VALUES (?, ?) RETURNING *;

-- name: GetProject :one
SELECT * FROM projects WHERE id = ? AND deleted_at IS NULL LIMIT 1;

-- name: ListProjects :many
SELECT * FROM projects WHERE deleted_at IS NULL ORDER BY id;

-- name: DeleteProject :exec
UPDATE projects SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL;
