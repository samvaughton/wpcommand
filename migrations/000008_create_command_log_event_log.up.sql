CREATE TABLE IF NOT EXISTS command_job_event_log(
    id SERIAL PRIMARY KEY,
    command_job_id INT NOT NULL,
    uuid TEXT UNIQUE NOT NULL,
    status TEXT NOT NULL,
    type TEXT NOT NULL,
    command TEXT DEFAULT '',
    output TEXT DEFAULT '',
    meta_data JSON DEFAULT '{}',
    created_at TIMESTAMP,
    executed_at TIMESTAMP,
    CONSTRAINT fk_command_job_id FOREIGN KEY (command_job_id) REFERENCES command_jobs (id)
);