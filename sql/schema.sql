-- users Table: Stores user information
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    permissions INTEGER NOT NULL DEFAULT 0,
    telegram_chat_id BIGINT UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- short_urls Table: Stores the mapping between short paths and original URLs
CREATE TABLE short_urls (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    short_path TEXT NOT NULL UNIQUE,
    original_url TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- url_clicks Table: Records each click for analytics
CREATE TABLE url_clicks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    short_url_id INTEGER NOT NULL,
    clicked_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    country_code TEXT,
    os_name TEXT,
    browser_name TEXT,
    raw_user_agent TEXT,
    FOREIGN KEY (short_url_id) REFERENCES short_urls(id)
);

-- telegram_auth_tokens Table: Stores temporary tokens for the Telegram account linking process
CREATE TABLE telegram_auth_tokens (
    token TEXT PRIMARY KEY,
    telegram_chat_id BIGINT NOT NULL,
    expires_at TIMESTAMP NOT NULL
);

-- Create a special 'anonymous' user for guest-created URLs
INSERT INTO users (username, password_hash, permissions) VALUES ('anonymous', '*', 0);
