package dependency

import (
	"ysf/dragonfly/tenantdb"
)

type defaultDep struct {
	tenantDB tenantdb.Service
}

func (d *defaultDep) TenantDB() tenantdb.Service {
	return d.tenantDB
}

func NewDefault(tenantDB tenantdb.Service) Service {
	return &defaultDep{
		tenantDB: tenantDB,
	}
}
