-- name: InsertClick :one
INSERT INTO url_clicks (
    short_url_id,
    country_code,
    os_name,
    browser_name,
    raw_user_agent
) VALUES (
    ?, ?, ?, ?, ?
)
RETURNING id;

-- name: CountClicksByShortURL :one
SELECT COUNT(*) FROM url_clicks
WHERE short_url_id = ?;

-- name: AggregateClicksByTimeRange :many
SELECT
    strftime('%Y-%m-%dT%H:00:00Z', clicked_at) as bucket_start,
    COUNT(*) as count
FROM url_clicks
WHERE
    short_url_id = ?
AND clicked_at BETWEEN ? AND ?
GROUP BY bucket_start
ORDER BY bucket_start;

-- name: AggregateClicksByCountry :many
SELECT
    country_code as agg_key,
    COUNT(*) as count
FROM url_clicks
WHERE
    short_url_id = ?
AND clicked_at BETWEEN ? AND ?
GROUP BY country_code
ORDER BY count DESC;

-- name: AggregateClicksByOS :many
SELECT
    os_name as agg_key,
    COUNT(*) as count
FROM url_clicks
WHERE
    short_url_id = ?
AND clicked_at BETWEEN ? AND ?
GROUP BY os_name
ORDER BY count DESC;

-- name: AggregateClicksByBrowser :many
SELECT
    browser_name as agg_key,
    COUNT(*) as count
FROM url_clicks
WHERE
    short_url_id = ?
AND clicked_at BETWEEN ? AND ?
GROUP BY browser_name
ORDER BY count DESC;