package migrate

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"
)

// CreateUsersTable1589341849 is struct to define a migration with ID 1589341849_create_users_table
type CreateUsersTable1589341849 struct{}

// ID return unique identifier for each migration. The prefix is unix time when this migration is created.
func (m CreateUsersTable1589341849) ID(ctx context.Context) string {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateUsersTable1589341849.ID")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	return fmt.Sprintf("%d_%s.sql", 1589341849, "create_users_table")
}

// SequenceNumber return current time when the migration is created,
// this useful to see the current status of the migration.
func (m CreateUsersTable1589341849) SequenceNumber(ctx context.Context) int {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateUsersTable1589341849.SequenceNumber")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	return 1589341849
}

// Up return sql migration for sync database
func (m CreateUsersTable1589341849) Up(ctx context.Context, tenantID string) (sql string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateUsersTable1589341849.Up")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	const query = `
CREATE TABLE IF NOT EXISTS %s.users (
	id UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
	username VARCHAR NOT NULL DEFAULT '',
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);
`
	sql = fmt.Sprintf(query, tenantID)
	return
}

// Down return sql migration for rollback database
func (m CreateUsersTable1589341849) Down(ctx context.Context, tenantID string) (sql string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateUsersTable1589341849.Down")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	const query = `DROP TABLE IF EXISTS %s.users;`
	sql = fmt.Sprintf(query, tenantID)
	return
}
