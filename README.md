# WPCOMMAND CLI

Kubernetes wordpress manager

# WPCOMMAND Server

## Migrations

We use https://github.com/golang-migrate/migrate

Example usage:

- Create a migration file: `bin/migrate.sh create -ext sql -dir ./db/migrations -seq create_foo_table`
- Run migrations: `bin/migrate.sh -path ./db/migrations up`

## Plugin/Theme Sync Architecture

A `Site` is assigned a `BlueprintSet` which is a `OneToMany` relationship of `ObjectBlueprints`. Each object blueprint 
will have its type (eg plugin), name, url and version amongst other fields. This set of object blueprints will then be
matched against the sites wordpress installation and synced.

Once an object has been created it will be versioned according to its version number provided. If a blueprint set
deployment goes wrong (ie updating an object to the next version and then apply the set) then the deployment set
can rollback the blueprint objects to the previous version.