CREATE TABLE IF NOT EXISTS object_blueprint_storage_relations(
   object_blueprint_id INT NOT NULL,
   object_blueprint_storage_id INT NOT NULL,
   PRIMARY KEY (object_blueprint_id, object_blueprint_storage_id),
   CONSTRAINT fk_object_blueprint_id FOREIGN KEY (object_blueprint_id) REFERENCES object_blueprints (id) ON DELETE CASCADE,
   CONSTRAINT fk_object_blueprint_storage_id FOREIGN KEY (object_blueprint_storage_id) REFERENCES object_blueprint_storage (id) ON DELETE CASCADE
);