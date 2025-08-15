-- users Table: Stores user information
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    permissions INTEGER NOT NULL DEFAULT 0,
    telegram_chat_id BIGINT UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_users_username
ON users(username)
WHERE deleted_at IS NULL;

-- short_urls Table: Stores the mapping between short paths and original URLs
CREATE TABLE IF NOT EXISTS short_urls (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    short_path TEXT NOT NULL,
    original_url TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_short_urls_short_path
ON short_urls(short_path)
WHERE deleted_at IS NULL;

-- url_clicks Table: Records each click for analytics
CREATE TABLE IF NOT EXISTS url_clicks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    short_url_id INTEGER NOT NULL,
    clicked_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    country_code TEXT,
    os_name TEXT,
    browser_name TEXT,
    raw_user_agent TEXT,
    ip_address TEXT,
    country TEXT,
    region_name TEXT,
    city TEXT,
    lat REAL,
    lon REAL,
    isp TEXT,
    as_info TEXT,
    is_processed BOOLEAN NOT NULL DEFAULT FALSE,
    is_success BOOLEAN NOT NULL DEFAULT TRUE,
    FOREIGN KEY (short_url_id) REFERENCES short_urls(id)
);

-- telegram_auth_tokens Table: Stores temporary tokens for the Telegram account linking process
CREATE TABLE IF NOT EXISTS telegram_auth_tokens (
    token TEXT PRIMARY KEY,
    telegram_chat_id BIGINT NOT NULL,
    expires_at TIMESTAMP NOT NULL
);
