CREATE TABLE IF NOT EXISTS accounts(
    id SERIAL PRIMARY KEY,
    uuid TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    key TEXT UNIQUE NOT NULL,
    enabled BOOLEAN,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

