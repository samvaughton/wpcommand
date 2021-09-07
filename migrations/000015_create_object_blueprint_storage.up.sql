CREATE TABLE IF NOT EXISTS object_blueprint_storage(
    id SERIAL PRIMARY KEY,
    uuid TEXT NOT NULL,
    hash TEXT NOT NULL,
    file BYTEA NOT NULL,
    created_at TIMESTAMP,
    UNIQUE(hash)
);