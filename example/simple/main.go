package main

import (
	"context"
	"fmt"
	"ysf/dragonfly/migration"
	"ysf/dragonfly/pkg/db"
	"ysf/dragonfly/tenantdb"
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

	// Then, for every other endpoint call, we can get the tenant's connection and do query based on tenant's data scope
	tenantConn, _ = tenant.GetTenant(ctx, tenantID)

	fmt.Println(tenantConn.TenantInfo())
	fmt.Println(tenantConn.ConnectionInfo())
}
