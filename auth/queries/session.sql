-- name: SessionCreate :one
INSERT INTO sessions (
    key,
    user_id, 
    created_at
    ) 
    VALUES (?, ?, ?)
    RETURNING id, key, user_id, created_at;

-- name: SessionRead :one
SELECT id, key, user_id ,created_at FROM sessions
WHERE key = ?;

-- name: SessionList :many
SELECT id, key, user_id, created_at FROM sessions
ORDER BY id ASC;

-- name: SessionTodayList :many
SELECT id, db_table, action, object_id, user_id, created_at, updated_at FROM logs
WHERE DATE(created_at) = DATE('now') AND db_table='session' And action='create';

-- name: SessionYesterdayList :many
SELECT id, db_table, action, object_id, user_id, created_at, updated_at FROM logs
WHERE DATE(created_at) = DATE('now', '-1 day') AND db_table='session' And action='create';

-- name: SessionPreviousWeeklyList :many
SELECT id, db_table, action, object_id, user_id, created_at, updated_at FROM logs
WHERE DATE(created_at) >= DATE('now', 'weekday 0') AND DATE(created_at) <= DATE('now') AND db_table='session' And action='create';

-- name: SessionWeeklyList :many
SELECT id, db_table, action, object_id, user_id, created_at, updated_at FROM logs
WHERE created_at >= DATE('now', 'weekday 0', '-7 days') AND created_at < DATE('now', 'weekday 0') AND db_table='session' And action='create';

-- name: SessionMonthlyList :many
SELECT id, db_table, action, object_id, user_id, created_at, updated_at FROM logs
WHERE strftime('%Y-%m', created_at) = strftime('%Y-%m', 'now') AND db_table='session' And action='create';

-- name: SessionPreviousMonthlyList :many
SELECT id, db_table, action, object_id, user_id, created_at, updated_at FROM logs
WHERE strftime('%Y-%m', created_at) = strftime('%Y-%m', 'now', '-1 month') AND db_table='session' And action='create';

-- name: SessionDelete :exec
DELETE FROM sessions WHERE key = ?;
