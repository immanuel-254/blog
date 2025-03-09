-- name: BlogCreate :one
INSERT INTO blogs (
    user_id, 
    title, 
    body, 
    created_at, 
    updated_at
    ) 
    VALUES (?, ?, ?, ?, ?)
    RETURNING *;

-- name: AssignBlogToCategory :exec
INSERT INTO category_blogs (blog_id, category_id, created_at, updated_at) VALUES (?, ?, ?, ?);

-- name: CategoryBlogDelete :exec
DELETE FROM category_blogs WHERE blog_id = ? and category_id = ?;

-- name: BlogList :many
SELECT id, user_id, title, body, created_at, updated_at FROM blogs
ORDER BY id ASC;

-- name: CategoryBlogList :many
SELECT blog_id, category_id, created_at, updated_at FROM category_blogs
ORDER BY blog_id ASC, category_id ASC;

-- name: BlogRead :one
SELECT id, user_id, title, body, created_at, updated_at FROM blogs
WHERE id = ?;

-- name: BlogCategoriesList :many
SELECT blog_id, category_id, created_at, updated_at FROM category_blogs
WHERE blog_id = ? ORDER BY blog_id ASC, category_id ASC;

-- name: BlogCommentsList :many
SELECT id, blog_id, user_id, body ,created_at, updated_at FROM comments
WHERE blog_id = ? ORDER BY blog_id ASC;

-- name: BlogUpdate :one
UPDATE blogs
SET 
    title = ?,
    body = ?,
    updated_at = ?
WHERE id = ?
RETURNING user_id, title, body, created_at, updated_at;

-- name: BlogDelete :exec
DELETE FROM blogs WHERE id = ?;
