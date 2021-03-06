package tenantdb

const (
	sqlCreateUuidExt = `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`

	// you can add postgres slaves as an object array contain {host, port} and username, password and db can use
	// same postgres master. You can also add redis cluster here
	sqlCreateConnectionTable = `
CREATE TABLE IF NOT EXISTS %s_connections (
	id UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
	postgres_master_debug BOOL NOT NULL DEFAULT true,
	postgres_master_host VARCHAR NOT NULL DEFAULT 'localhost',
	postgres_master_port INT NOT NULL DEFAULT 5432,
	postgres_master_username VARCHAR NOT NULL DEFAULT '',
	postgres_master_password VARCHAR NOT NULL DEFAULT '',
	postgres_master_database VARCHAR NOT NULL DEFAULT '',
	postgres_master_pool_size INT NOT NULL DEFAULT 10,
	postgres_master_read_timeout INT NOT NULL DEFAULT 1000,
	postgres_master_write_timeout INT NOT NULL DEFAULT 1000
);
`

	sqlCreateTenantsTable = `
CREATE TABLE IF NOT EXISTS %s_tenants (
	id VARCHAR NOT NULL PRIMARY KEY,
	name VARCHAR NOT NULL DEFAULT '',
	connection_id UUID NOT NULL CONSTRAINT %s_tenants_connection_id_foreign REFERENCES %s_connections(id) ON UPDATE CASCADE ON DELETE RESTRICT,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_idx_tenants_id ON %s_tenants (id);

COMMENT ON COLUMN %s_tenants.id IS 'useful for creating postgres schema and table prefix, so it will not be conflict with another tenant in same db';
`

	sqlGetTenantByID = `SELECT * FROM %s_tenants WHERE id = ?0 LIMIT 1;`

	sqlGetTenants = `SELECT * FROM %s_tenants;`

	sqlGetConnectionById = `SELECT * FROM %s_connections WHERE id = ?0 LIMIT 1;`

	sqlGetConnections = `SELECT * FROM %s_connections;`

	sqlCreateTenant = `INSERT INTO %s_tenants (id, name, connection_id) VALUES (?0, ?1, ?2) RETURNING *;`
)
