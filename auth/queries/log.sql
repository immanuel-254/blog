-- name: LogCreate :exec
INSERT INTO logs (
    db_table, 
    action,
    object_id, 
    user_id, 
    created_at, 
    updated_at
    ) 
    VALUES (?, ?, ?, ?, ?, ?);

-- name: LogList :many
SELECT id, db_table, action, object_id, user_id, created_at, updated_at FROM logs
ORDER BY id ASC;

-- name: LogTodayList :many
SELECT id, db_table, action, object_id, user_id, created_at, updated_at FROM logs
WHERE DATE(created_at) = DATE('now');

-- name: LogYesterdayList :many
SELECT id, db_table, action, object_id, user_id, created_at, updated_at FROM logs
WHERE DATE(created_at) = DATE('now', '-1 day');

-- name: LogPreviousWeeklyList :many
SELECT id, db_table, action, object_id, user_id, created_at, updated_at FROM logs
WHERE DATE(created_at) >= DATE('now', 'weekday 0') AND DATE(created_at) <= DATE('now');

-- name: LogWeeklyList :many
SELECT id, db_table, action, object_id, user_id, created_at, updated_at FROM logs
WHERE created_at >= DATE('now', 'weekday 0', '-7 days') AND created_at < DATE('now', 'weekday 0');

-- name: LogMonthlyList :many
SELECT id, db_table, action, object_id, user_id, created_at, updated_at FROM logs
WHERE strftime('%Y-%m', created_at) = strftime('%Y-%m', 'now');

-- name: LogPreviousMonthlyList :many
SELECT id, db_table, action, object_id, user_id, created_at, updated_at FROM logs
WHERE strftime('%Y-%m', created_at) = strftime('%Y-%m', 'now', '-1 month');