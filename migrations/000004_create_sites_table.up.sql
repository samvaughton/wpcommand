CREATE TABLE IF NOT EXISTS sites(
    id SERIAL PRIMARY KEY,
    account_id INT NOT NULL,
    uuid TEXT UNIQUE NOT NULL,
    key TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL,
    label_selector TEXT DEFAULT '',
    namespace TEXT DEFAULT '',
    enabled BOOLEAN,
    site_email TEXT DEFAULT '',
    site_username TEXT DEFAULT '',
    site_password TEXT DEFAULT '',
    site_config JSON DEFAULT '{}',
    wp_cached_data JSON DEFAULT '{}',
    test_mode BOOL DEFAULT FALSE NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT fk_account_id FOREIGN KEY (account_id) REFERENCES accounts (id)
);