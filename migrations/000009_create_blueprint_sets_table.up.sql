CREATE TABLE IF NOT EXISTS blueprint_sets(
    id SERIAL PRIMARY KEY,
    uuid TEXT UNIQUE NOT NULL,
    name TEXT UNIQUE NOT NULL,
    enabled BOOLEAN,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);