-- name: CreateShortURL :one
INSERT INTO short_urls (short_path, original_url, user_id)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetShortURLByPath :one
SELECT *
FROM short_urls
WHERE short_path = ?;

-- name: GetShortURLByID :one
SELECT *
FROM short_urls
WHERE id = ?;

-- TODO: short delete

-- name: DeleteShortURL :exec
DELETE FROM short_urls
WHERE id = ?;

-- name: ListShortURLsByUserID :many
SELECT
    sqlc.embed(su),
    (SELECT COUNT(*) FROM url_clicks uc WHERE uc.short_url_id = su.id) AS total_clicks
FROM short_urls su
WHERE su.user_id = ?
ORDER BY su.created_at DESC;

-- name: ListAllShortURLs :many
SELECT
    sqlc.embed(su),
    (SELECT COUNT(*) FROM url_clicks uc WHERE uc.short_url_id = su.id) AS total_clicks
FROM short_urls su
JOIN users u ON su.user_id = u.id
ORDER BY su.created_at DESC;
