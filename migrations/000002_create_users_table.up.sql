CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    uuid TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    first_name TEXT DEFAULT '',
    last_name TEXT DEFAULT '',
    password TEXT DEFAULT '',
    enabled BOOLEAN,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);


