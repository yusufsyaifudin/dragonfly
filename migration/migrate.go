package migration

import "context"

type Migrate interface {
	// ID return unique identifier for each migration
	ID(ctx context.Context) string

	// Up return sql migration for sync database
	Up(ctx context.Context, tenantID string) (sql string, err error)

	// Down return sql migration for rollback database
	Down(ctx context.Context, tenantID string) (sql string, err error)
}
