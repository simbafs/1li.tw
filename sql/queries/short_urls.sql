-- name: CreateShortURL :one
INSERT INTO short_urls (short_path, original_url, user_id)
VALUES (?, ?, ?)
RETURNING id, short_path, original_url;

-- name: GetShortURLByPath :one
SELECT id, original_url, user_id
FROM short_urls
WHERE short_path = ?;

-- name: GetShortURLByID :one
SELECT id, short_path, original_url, user_id
FROM short_urls
WHERE id = ?;

-- name: DeleteShortURL :exec
DELETE FROM short_urls
WHERE id = ?;

-- name: ListShortURLsByUserID :many
SELECT
    su.id,
    su.short_path,
    su.original_url,
    su.created_at,
    (SELECT COUNT(*) FROM url_clicks uc WHERE uc.short_url_id = su.id) AS total_clicks
FROM short_urls su
WHERE su.user_id = ?
ORDER BY su.created_at DESC;

-- name: ListAllShortURLs :many
SELECT
    su.id,
    su.short_path,
    su.original_url,
    su.created_at,
    su.user_id,
    u.username as owner_username,
    (SELECT COUNT(*) FROM url_clicks uc WHERE uc.short_url_id = su.id) AS total_clicks
FROM short_urls su
JOIN users u ON su.user_id = u.id
ORDER BY su.created_at DESC;