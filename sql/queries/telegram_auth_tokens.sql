-- sql/queries/telegram_auth_tokens.sql

-- name: CreateTelegramAuthToken :exec
INSERT INTO telegram_auth_tokens (token, telegram_chat_id, expires_at)
VALUES (?, ?, ?);

-- name: GetTelegramAuthToken :one
SELECT token, telegram_chat_id, expires_at FROM telegram_auth_tokens
WHERE token = ? AND expires_at > CURRENT_TIMESTAMP LIMIT 1;

-- name: DeleteTelegramAuthToken :exec
DELETE FROM telegram_auth_tokens
WHERE token = ?;
