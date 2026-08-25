-- name: BranchCreate :exec
INSERT INTO branches (id, project_id, name, commit_id) VALUES (?, ?, ?, ?);

-- name: BranchUpdate :exec
UPDATE branches SET name = ?, commit_id = ?, protected = ? WHERE id = ?;

-- name: BranchList :many
SELECT id, project_id, name, protected, commit_id, updated_at, created_at
FROM branches
WHERE project_id = :project_id
  AND (updated_at < sqlc.arg(last_updated_at) is null or updated_at < sqlc.arg(last_updated_at) OR (updated_at = sqlc.arg(last_updated_at) AND id < sqlc.arg(last_id)))
ORDER BY updated_at DESC, id DESC
LIMIT :limit;