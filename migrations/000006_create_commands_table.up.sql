CREATE TABLE IF NOT EXISTS commands(
    id SERIAL PRIMARY KEY,
    account_id INT,
    uuid TEXT UNIQUE NOT NULL,
    type TEXT NOT NULL,
    key TEXT NOT NULL,
    description TEXT,
    http_method TEXT,
    http_url TEXT,
    http_headers JSON,
    http_body TEXT,
    created_at TIMESTAMP,
    CONSTRAINT fk_account_id FOREIGN KEY (account_id) REFERENCES accounts (id)
);