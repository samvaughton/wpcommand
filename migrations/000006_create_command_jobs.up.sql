CREATE TABLE IF NOT EXISTS command_jobs(
    id SERIAL PRIMARY KEY,
    site_id INT NOT NULL,
    uuid TEXT UNIQUE NOT NULL,
    status TEXT NOT NULL,
    key TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP,
    CONSTRAINT fk_site_id FOREIGN KEY (site_id) REFERENCES sites (id)
);