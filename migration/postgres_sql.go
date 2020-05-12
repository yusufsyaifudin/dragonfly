package migration

const (
	// recording migration
	sqlCreateMigrationTable = `
	CREATE TABLE IF NOT EXISTS %s.%s_migrations (
		id VARCHAR NOT NULL PRIMARY KEY,
		sequence_number INT NOT NULL,
		applied_at TIMESTAMP WITH TIMEZONE NOT NULL DEFAULT now()
	);
`

	// select all applied migration
	sqlGetMigrationData = `SELECT * FROM %s.%s_migrations ORDER BY sequence_number ASC;`

	sqlInsertMigrationData = `
	INSERT INTO %s.%s_migrations (id, sequence_number, applied_at)
	VALUES (?0, ?1, now()) 
	RETURNING *;
`
)
