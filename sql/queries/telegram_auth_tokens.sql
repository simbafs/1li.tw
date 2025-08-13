-- name: CreateTelegramAuthToken :one
INSERT INTO telegram_auth_tokens (token, telegram_chat_id, expires_at)
VALUES (?, ?, ?)
RETURNING token;

-- name: GetTelegramAuthToken :one
SELECT token, telegram_chat_id, expires_at
FROM telegram_auth_tokens
WHERE token = ?;

-- name: DeleteTelegramAuthToken :exec
DELETE FROM telegram_auth_tokens
WHERE token = ?;
