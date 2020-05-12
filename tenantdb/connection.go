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

type defaultConnection struct {
	tenantInfo *model.Tenant
	connInfo   *model.Connection
	sql        db.SQL
}

func (d defaultConnection) TenantInfo() *model.Tenant {
	return d.tenantInfo
}

func (d defaultConnection) ConnectionInfo() *model.Connection {
	return d.connInfo
}

func (d defaultConnection) SQL() db.SQL {
	return d.sql
}

func (d defaultConnection) setSQL(sql db.SQL) {
	d.sql = sql
}

func (d defaultConnection) Redis() {
	panic("implement me")
}

func NewConnection(tenantInfo *model.Tenant, connInfo *model.Connection) *defaultConnection {
	if tenantInfo == nil {
		panic("cannot create connection, tenant info is nil")
	}

	if connInfo == nil {
		panic("cannot create connection, connection info is nil")
	}

	return &defaultConnection{
		tenantInfo: tenantInfo,
		connInfo:   connInfo,
	}
}
