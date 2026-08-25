-- name: BranchCreate :exec
INSERT INTO branches (id, project_id, name, commit_id) VALUES (?, ?, ?, ?);

-- name: BranchUpdate :exec
UPDATE branches SET name = ?, key = ?, commit_id = ?, protected = ? WHERE id = :id;

-- name: BranchList :many
SELECT *
FROM branches
WHERE project_id = :project_id
  AND (sqlc.arg(last_updated_at) is null or updated_at < sqlc.arg(last_updated_at) OR (updated_at = sqlc.arg(last_updated_at) AND id < sqlc.arg(last_id)))
ORDER BY updated_at DESC, id DESC
LIMIT :limit;

-- name: BranchGet :one
SELECT *
FROM branches
WHERE project_id = :project_id AND id = :id
LIMIT 1;