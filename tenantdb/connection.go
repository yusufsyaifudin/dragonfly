package tenantdb

import (
	"ysf/dragonfly/pkg/db"
	"ysf/dragonfly/tenantdb/model"
)

// Connection real connection to specific tenant
type Connection interface {
	TenantInfo() *model.Tenant
	ConnectionInfo() *model.Connection
	SQL() db.SQL
	Redis()
}
