-- name: CreateUser :one
INSERT INTO users (username, password_hash, permissions)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = ?;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = ?;

-- name: GetUserByTelegramID :one
SELECT *
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
