package tenantdb

import (
	"context"
	"ysf/dragonfly/migration"
	"ysf/dragonfly/tenantdb/model"
)

type Service interface {
	// CreateTenant will return connector (immigration and connection) so it can be migrated and created
	CreateTenant(ctx context.Context, tenantId, tenantName string) (connection Connection, err error)
	GetTenant(ctx context.Context, tenantId string) (connection Connection, err error)
	GetTenants(ctx context.Context) (tenants model.Tenants, err error)

	GetTenantImmigration(ctx context.Context, tenantId string) (immigration migration.Immigration, err error)
}
