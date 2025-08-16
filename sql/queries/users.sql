-- name: CreateUser :one
INSERT INTO users (username, password_hash, permissions)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = ? AND deleted_at IS NULL;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = ? AND deleted_at IS NULL;

-- name: GetUserByTelegramID :one
SELECT *
FROM users
WHERE telegram_chat_id = ? AND deleted_at IS NULL;

-- name: UpdateUserTelegramID :exec
UPDATE users
SET telegram_chat_id = ?
WHERE id = ?;

-- name: UpdateUserPermissions :exec
UPDATE users
SET permissions = ?
WHERE id = ?;

-- name: DeleteUser :exec
UPDATE users
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: ListUsers :many
SELECT * 
FROM users
WHERE deleted_at IS NULL
ORDER BY permissions DESC, id ASC;
