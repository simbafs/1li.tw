-- name: CreateUser :one
INSERT INTO users (
    username,
    password_hash,
    permissions
) VALUES (
    ?, ?, ?
)
RETURNING id;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = ?;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = ?;

-- name: UpdateUserTelegramID :exec
UPDATE users
SET telegram_chat_id = ?
WHERE id = ?;

-- name: UpdateUserPermissions :exec
UPDATE users
SET permissions = ?
WHERE id = ?;