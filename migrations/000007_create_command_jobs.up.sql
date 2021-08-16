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
    CONSTRAINT fk_command_id FOREIGN KEY (command_id) REFERENCES commands (id),
    CONSTRAINT fk_site_id FOREIGN KEY (site_id) REFERENCES sites (id),
    CONSTRAINT fk_run_by_user_id FOREIGN KEY (run_by_user_id) REFERENCES users (id)
);