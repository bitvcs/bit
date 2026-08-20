-- name: CreateProject :one
INSERT INTO projects (org_id, slug, name, description) VALUES (?, ?, ?, ?) RETURNING *;

-- name: GetProject :one
SELECT * FROM projects WHERE id = ? AND deleted = false LIMIT 1;

-- name: ListProjectsByOrgId :many
SELECT * FROM projects WHERE org_id = ? AND deleted = false ORDER BY id;

-- name: DeleteProject :exec
UPDATE projects SET deleted = true, deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted = false;
