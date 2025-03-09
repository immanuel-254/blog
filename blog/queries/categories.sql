-- name: CategoryCreate :one
INSERT INTO categories (
    user_id, 
    name, 
    created_at, 
    updated_at
    ) 
    VALUES (?, ?, ?, ?)
    RETURNING *;

-- name: CategoryList :many
SELECT id, user_id, name, created_at, updated_at FROM categories
ORDER BY id ASC;

-- name: CategoryRead :one
SELECT id, user_id, name, created_at, updated_at FROM categories
WHERE id = ?;

-- name: CategoryUpdate :one
UPDATE categories
SET 
    name = ?,
    updated_at = ?
WHERE id = ?
RETURNING id, user_id, name, created_at, updated_at;

-- name: CategoryDelete :exec
DELETE FROM categories WHERE id = ?