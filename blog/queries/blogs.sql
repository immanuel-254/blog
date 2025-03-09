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
SELECT 
    b.id AS blog_id,
    b.title AS blog_title,
    b.body AS blog_body,
    b.publish AS blog_publish,
    b.created_at AS blog_created_at,
    b.updated_at AS blog_updated_at,
    p.user_id AS blog_auth_id,
    p.username AS user_name,

    COALESCE(GROUP_CONCAT(DISTINCT 
        CONCAT('{"id":', cat.id, ',"name":"', cat.name, '"}')    
    ), '[]') AS categories,
    COALESCE(GROUP_CONCAT(DISTINCT 
        CONCAT('{"id":', c.id, ',"user_id":', c.user_id, ',"body":"', c.body, '"}')
    ), '[]') AS comments
    
FROM blogs b
LEFT JOIN profiles p ON b.user_id = p.user_id
LEFT JOIN comments c ON b.id = c.blog_id
LEFT JOIN category_blogs cb ON b.id = cb.blog_id
LEFT JOIN categories cat ON cb.category_id = cat.id
WHERE b.id = ?; 

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
