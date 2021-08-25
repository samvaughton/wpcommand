CREATE TABLE IF NOT EXISTS commands(
    id SERIAL PRIMARY KEY,
    account_id INT,
    site_id INT,
    uuid TEXT UNIQUE NOT NULL,
    type TEXT NOT NULL,
    key TEXT NOT NULL,
    description TEXT,
    http_method TEXT,
    http_url TEXT,
    http_headers JSON,
    http_body TEXT,
    public BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP
);