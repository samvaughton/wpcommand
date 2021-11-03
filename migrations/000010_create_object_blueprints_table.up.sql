CREATE TABLE IF NOT EXISTS object_blueprints(
     id SERIAL PRIMARY KEY,
     revision_id INT,
     blueprint_set_id INT NOT NULL,
     set_order INT NOT NULL,
     uuid TEXT NOT NULL,
     type TEXT NOT NULL,
     name TEXT NOT NULL,
     exact_name TEXT NOT NULL,
     version TEXT NOT NULL,
     original_object_url TEXT NOT NULL,
     active BOOLEAN DEFAULT FALSE,
     created_at TIMESTAMP,
     updated_at TIMESTAMP,
     CONSTRAINT fk_blueprint_set_id FOREIGN KEY (blueprint_set_id) REFERENCES blueprint_sets (id) ON DELETE CASCADE,
     CONSTRAINT object_blueprints_uuid_revision_id_uindex UNIQUE (uuid, revision_id),
     CONSTRAINT object_blueprints_uuid_version_uindex UNIQUE (uuid, version)
);