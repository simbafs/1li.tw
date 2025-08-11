-- sql/queries/users.sql

-- name: CreateUser :one
INSERT INTO users (username, password_hash, permissions, telegram_chat_id)
VALUES (?, ?, ?, ?)
RETURNING id, username, permissions, created_at;

-- name: GetUserByUsername :one
SELECT id, username, password_hash, permissions, telegram_chat_id, created_at FROM users
WHERE username = ? LIMIT 1;

-- name: GetUserByID :one
SELECT id, username, password_hash, permissions, telegram_chat_id, created_at FROM users
WHERE id = ? LIMIT 1;

-- name: UpdateUserTelegramID :exec
UPDATE users
SET telegram_chat_id = ?
WHERE id = ?;

-- name: UpdateUserRole :exec
UPDATE users
SET permissions = ?
WHERE id = ?;
