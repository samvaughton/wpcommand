CREATE TABLE IF NOT EXISTS sites_blueprint_sets(
     site_id INT NOT NULL,
     blueprint_set_id INT NOT NULL,
    PRIMARY KEY (site_id, blueprint_set_id),
    CONSTRAINT fk_site_id FOREIGN KEY (site_id) REFERENCES sites (id),
    CONSTRAINT fk_blueprint_set_id FOREIGN KEY (blueprint_set_id) REFERENCES blueprint_sets (id)
);