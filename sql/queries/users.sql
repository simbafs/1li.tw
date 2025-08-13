-- name: CreateUser :one
INSERT INTO users (username, password_hash, permissions)
VALUES (?, ?, ?)
RETURNING id, username;

-- name: GetUserByUsername :one
SELECT id, username, password_hash, permissions, telegram_chat_id, created_at
FROM users
WHERE username = ?;

-- name: GetUserByID :one
SELECT id, username, password_hash, permissions, telegram_chat_id, created_at
FROM users
WHERE id = ?;

-- name: GetUserByTelegramID :one
SELECT id, username, password_hash, permissions, telegram_chat_id, created_at
FROM users
WHERE telegram_chat_id = ?;

-- name: UpdateUserTelegramID :exec
UPDATE users
SET telegram_chat_id = ?
WHERE id = ?;

-- name: UpdateUserPermissions :exec
UPDATE users
SET permissions = ?
WHERE id = ?;