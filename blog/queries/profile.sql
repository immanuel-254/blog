-- name: ProfileCreate :one
INSERT INTO profiles (
    user_id, 
    username,
    image,
    bio, 
    created_at, 
    updated_at
    ) 
    VALUES (?, ?, ?, ?, ?, ?)
    RETURNING *;

-- name: ProfileList :many
SELECT id, user_id, username, image, bio, created_at, updated_at FROM profiles
ORDER BY id ASC;

-- name: ProfileRead :one
SELECT id, user_id, username, image, bio, created_at, updated_at FROM profiles
WHERE id = ?;

-- name: ProfileUpdate :one
UPDATE profiles
SET 
    username = ?,
    image = ?,
    bio = ?,
    created_at = ?,
    updated_at = ?
WHERE id = ?
RETURNING id, user_id, username, image, bio, created_at, updated_at;

-- name: ProfileDelete :exec
DELETE FROM profiles WHERE id = ?
