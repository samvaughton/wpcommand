CREATE TABLE IF NOT EXISTS settings(
   id SERIAL PRIMARY KEY,
   uuid TEXT UNIQUE NOT NULL,
   name TEXT NOT NULL,
   value TEXT,
   created_at TIMESTAMP,
   updated_at TIMESTAMP
);

