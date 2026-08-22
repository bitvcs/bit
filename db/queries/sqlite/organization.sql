-- name: CreateOrganization :one
INSERT INTO organizations (id, name, slug) VALUES (?, ?, ?) RETURNING *;

-- name: GetOrganization :one
SELECT * FROM organizations WHERE id = ? AND deleted = false LIMIT 1;

-- name: ListOrganizations :many
SELECT * FROM organizations WHERE deleted = false ORDER BY id;

-- name: DeleteOrganization :exec
UPDATE organizations SET deleted = true, deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted = false;
