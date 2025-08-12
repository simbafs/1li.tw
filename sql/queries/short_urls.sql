-- name: CreateShortURL :one
INSERT INTO short_urls (
    short_path,
    original_url,
    user_id
) VALUES (
    ?, ?, ?
)
RETURNING id;

-- name: GetShortURLByPath :one
SELECT * FROM short_urls
WHERE short_path = ?;

-- name: DeleteShortURL :exec
DELETE FROM short_urls
WHERE id = ?;

-- name: ListShortURLsByUser :many
SELECT * FROM short_urls
WHERE user_id = ?
ORDER BY created_at DESC;