CREATE TABLE IF NOT EXISTS command_jobs(
    id SERIAL PRIMARY KEY,
    command_id INT NOT NULL,
    run_by_user_id INT,
    site_id INT NOT NULL,
    uuid TEXT UNIQUE NOT NULL,
    status TEXT NOT NULL,
    key TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP,
    config JSON DEFAULT '{}'::json NOT NULL,
    result_data JSON DEFAULT '{}'::json NOT NULL,
    CONSTRAINT fk_command_id FOREIGN KEY (command_id) REFERENCES commands (id),
    CONSTRAINT fk_site_id FOREIGN KEY (site_id) REFERENCES sites (id)
);