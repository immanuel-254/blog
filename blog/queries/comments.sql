-- name: CommentCreate :one
INSERT INTO comments (
    user_id, 
    body, 
    created_at, 
    updated_at
    ) 
    VALUES (?, ?, ?, ?)
    RETURNING *;

-- name: CommentList :many
SELECT id, user_id, body, created_at, updated_at FROM comments
ORDER BY id ASC;

-- name: CommentRead :one
SELECT id, user_id, body, created_at, updated_at FROM comments
WHERE id = ?;

-- name: CommentUpdate :one
UPDATE comments
SET 
    body = ?,
    updated_at = ?
WHERE id = ?
RETURNING id, user_id, body, created_at, updated_at;

-- name: CommentDelete :exec
DELETE FROM comments WHERE id = ?
