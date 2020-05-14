package dependency

import (
	"ysf/dragonfly/tenantdb"
)

type Service interface {
	TenantDB() tenantdb.Service
}
