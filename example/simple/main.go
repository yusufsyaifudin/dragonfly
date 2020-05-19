package main

import (
	"context"
	"fmt"
	"ysf/dragonfly/pkg/db"
	"ysf/dragonfly/tenantdb"
	"ysf/dragonfly/tenantdb/migration"
)

// note that this program is just give you a vision about how using the API.
// You must create real connection to database to make it works
func main() {
	ctx := context.Background()
	const tenantID = "tenant1"

	conn := db.NewNoop()

	// migration must per each scope,
	m := []migration.Migrate{
		nil,
	}

	// This should be initiated once in application start
	tenant, _ := tenantdb.Postgres(conn, m)

	// When calling API create new tenant (for example register endpoint), we must create new tenant and do a migration
	// to prepare the box or repository for tenant's data
	tenantConn, _ := tenant.CreateTenant(ctx, tenantID, tenantID)
	tenantMigration, _ := tenant.GetTenantImmigration(ctx, tenantID)
	_ = tenantMigration.Sync(ctx)

	// example to handle different case, for example we want to defer migration on specific tenant id,
	// so we must use different query to handle it
	query := `SELECT * FROM product_id = ?0 LIMIT 1;`
	if tenantMigration.Status(ctx).LastSequenceNumber <= 10 {
		query = `SELECT * FROM product_code = ?0 LIMIT 1;`
	}

	_ = tenantConn.SQL().Reader().ExecContext(ctx, query)

	// Then, for every other endpoint call, we can get the tenant's connection and do query based on tenant's data scope
	tenantConn, _ = tenant.GetTenant(ctx, tenantID)

	fmt.Println(tenantConn.TenantInfo())
	fmt.Println(tenantConn.ConnectionInfo())
}
