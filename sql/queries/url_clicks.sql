-- name: CreateURLClick :one
INSERT INTO url_clicks (short_url_id, country_code, os_name, browser_name, raw_user_agent, ip_address)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING id;

-- name: CountClicksByShortURLID :one
SELECT COUNT(*)
FROM url_clicks
WHERE short_url_id = ?;

-- name: GetClickStatsByTime :many
SELECT
    strftime('%Y-%m-%dT%H:00:00Z', clicked_at) as time_bucket,
    COUNT(*) as count
FROM url_clicks
WHERE short_url_id = ? AND clicked_at >= sqlc.arg('from') AND clicked_at <= sqlc.arg('to')
GROUP BY time_bucket
ORDER BY time_bucket;

-- name: GetClickStatsByCountry :many
SELECT
    country_code,
    COUNT(*) as count
FROM url_clicks
WHERE short_url_id = ? AND clicked_at >= sqlc.arg('from') AND clicked_at <= sqlc.arg('to')
GROUP BY country_code
ORDER BY count DESC;

-- name: GetClickStatsByOS :many
SELECT
    os_name,
    COUNT(*) as count
FROM url_clicks
WHERE short_url_id = ? AND clicked_at >= sqlc.arg('from') AND clicked_at <= sqlc.arg('to')
GROUP BY os_name
ORDER BY count DESC;

-- name: GetClickStatsByBrowser :many
SELECT
    browser_name,
    COUNT(*) as count
FROM url_clicks
WHERE short_url_id = ? AND clicked_at >= sqlc.arg('from') AND clicked_at <= sqlc.arg('to')
GROUP BY browser_name
ORDER BY count DESC;

-- name: GetUnprocessedClicks :many
SELECT id, ip_address
FROM url_clicks
WHERE is_processed = FALSE AND ip_address IS NOT NULL AND ip_address != ''
LIMIT ?;

-- name: UpdateClickGeoInfo :exec
UPDATE url_clicks
SET
    is_success = ?,
    country = ?,
    region_name = ?,
    city = ?,
    lat = ?,
    lon = ?,
    isp = ?,
    as_info = ?,
    is_processed = TRUE
WHERE ip_address = ?;
