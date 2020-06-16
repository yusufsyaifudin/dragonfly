package tenantdb

import (
	"context"

	"github.com/yusufsyaifudin/dragonfly/tenantdb/migration"
	"github.com/yusufsyaifudin/dragonfly/tenantdb/model"
)

type Service interface {
	// CreateTenant will return connector (immigration and connection) so it can be migrated and created
	CreateTenant(ctx context.Context, tenantId, tenantName string) (tenant *model.Tenant, connection Connection, err error)
	GetTenant(ctx context.Context, tenantId string) (tenant *model.Tenant, connection Connection, err error)
	GetTenants(ctx context.Context) (tenants model.Tenants, err error)

	GetTenantImmigration(ctx context.Context, tenantId string) (immigration migration.Immigration, err error)

	// Close will close all established connection foreach tenant connection
	Close() error
}
