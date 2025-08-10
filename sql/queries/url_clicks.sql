-- sql/queries/url_clicks.sql

-- name: InsertClick :one
INSERT INTO url_clicks (short_url_id, country_code, user_agent)
VALUES (?, ?, ?)
RETURNING id;

-- name: GetClicksByShortURLID :many
SELECT id, short_url_id, clicked_at, country_code, user_agent FROM url_clicks
WHERE short_url_id = ?
ORDER BY clicked_at DESC;

-- name: AggregateClicksByTime :many
SELECT
  strftime('%Y-%m-%dT%H:00:00Z', clicked_at) AS bucket_start,
  COUNT(*) AS count
FROM url_clicks
WHERE short_url_id = ? AND clicked_at BETWEEN ? AND ?
GROUP BY bucket_start
ORDER BY bucket_start;

-- name: AggregateClicksByCountry :many
SELECT country_code AS country_key, COUNT(*) AS count FROM url_clicks
WHERE short_url_id = ? AND country_code IS NOT NULL AND clicked_at BETWEEN ? AND ?
GROUP BY country_key
ORDER BY count DESC;

-- name: AggregateClicksByOS :many
SELECT
  CASE
    WHEN user_agent LIKE '%Android%' THEN 'Android'
    WHEN user_agent LIKE '%iPhone%' OR user_agent LIKE '%iPad%' THEN 'iOS'
    WHEN user_agent LIKE '%Windows%' THEN 'Windows'
    WHEN user_agent LIKE '%Mac OS%' THEN 'macOS'
    WHEN user_agent LIKE '%Linux%' THEN 'Linux'
    ELSE 'Other'
  END AS os_name,
  COUNT(*) AS count
FROM url_clicks
WHERE short_url_id = ? AND user_agent IS NOT NULL AND clicked_at BETWEEN ? AND ?
GROUP BY os_name
ORDER BY count DESC;

-- name: AggregateClicksByBrowser :many
SELECT
  CASE
    WHEN user_agent LIKE '%Chrome%' AND user_agent NOT LIKE '%Chromium%' THEN 'Chrome'
    WHEN user_agent LIKE '%Firefox%' THEN 'Firefox'
    WHEN user_agent LIKE '%Safari%' AND user_agent NOT LIKE '%Chrome%' THEN 'Safari'
    WHEN user_agent LIKE '%Edge%' THEN 'Edge'
    WHEN user_agent LIKE '%Opera%' THEN 'Opera'
    ELSE 'Other'
  END AS browser_name,
  COUNT(*) AS count
FROM url_clicks
WHERE short_url_id = ? AND user_agent IS NOT NULL AND clicked_at BETWEEN ? AND ?
GROUP BY browser_name
ORDER BY count DESC;
