package migration

const (
	// create schema
	sqlCreatePostgresSchema = `
	CREATE SCHEMA IF NOT EXISTS %s;
`

	// recording migration
	sqlCreateMigrationTable = `
	CREATE TABLE IF NOT EXISTS %s.migrations (
		id VARCHAR NOT NULL PRIMARY KEY,
		sequence_number INT NOT NULL,
		applied_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
	);

	CREATE UNIQUE INDEX IF NOT EXISTS unique_%s_migrations_seq_num ON %s.migrations (sequence_number);
`

	// select all applied migration
	sqlGetMigrationData = `SELECT * FROM %s.migrations ORDER BY sequence_number ASC;`

	sqlInsertMigrationData = `
	INSERT INTO %s.migrations (id, sequence_number, applied_at)
	VALUES (?0, ?1, now()) 
	RETURNING *;
`
)
