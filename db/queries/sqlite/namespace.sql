-- name: CreateNamespace :one
INSERT INTO namespaces (name) VALUES (?) RETURNING *;

-- name: GetNamespace :one
SELECT * FROM namespaces WHERE id = ? AND deleted_at IS NULL LIMIT 1;

-- name: ListNamespaces :many
SELECT * FROM namespaces WHERE deleted_at IS NULL ORDER BY id;

-- name: DeleteNamespace :exec
UPDATE namespaces SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL;
