CREATE TABLE IF NOT EXISTS blueprint_sets(
    id SERIAL PRIMARY KEY,
    account_id INT NOT NULL,
    uuid TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    enabled BOOLEAN,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT fk_account_id FOREIGN KEY (account_id) REFERENCES accounts (id)
);