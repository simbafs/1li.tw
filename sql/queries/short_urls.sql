-- sql/queries/short_urls.sql

-- name: CreateShortURL :one
INSERT INTO short_urls (short_path, original_url, user_id)
VALUES (?, ?, ?)
RETURNING id, short_path, original_url, user_id, created_at;

-- name: GetShortURLByPath :one
SELECT id, short_path, original_url, user_id, created_at FROM short_urls
WHERE short_path = ? LIMIT 1;

-- name: DeleteShortURL :exec
DELETE FROM short_urls
WHERE id = ?;

-- name: ListShortURLsByUser :many
SELECT id, short_path, original_url, user_id, created_at FROM short_urls
WHERE user_id = ?
ORDER BY created_at DESC;

-- name: ListAllShortURLs :many
SELECT su.id, su.short_path, su.original_url, su.user_id, su.created_at, u.username AS owner_username
FROM short_urls su
JOIN users u ON su.user_id = u.id
ORDER BY su.created_at DESC;

-- name: CountClicksByShortURL :one
SELECT COUNT(*) FROM url_clicks
WHERE short_url_id = ?;
