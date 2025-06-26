CREATE TABLE IF NOT EXISTS tokens (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    refresh_token TEXT NOT NULL,
    expires_at    TEXT NOT NULL -- ISO8601 timestamp
);