package tenantdb

import (
	"fmt"
	"ysf/dragonfly/pkg/db"
	"ysf/dragonfly/tenantdb/model"
)

type defaultConnection struct {
	connInfo *model.Connection
	sql      db.SQL
}

func (d *defaultConnection) ConnectionInfo() *model.Connection {
	return d.connInfo
}

func (d *defaultConnection) SQL() db.SQL {
	return d.sql
}

func (d *defaultConnection) setSQL(sql db.SQL) {
	d.sql = sql
}

func (d *defaultConnection) Redis() {
	panic("implement me")
}

func newConnection(connInfo *model.Connection) (connector *defaultConnection, err error) {
	connector = &defaultConnection{
		connInfo: connInfo,
	}

	if connInfo == nil {
		err = fmt.Errorf("cannot create connection, connection info is nil")
		return
	}

	return
}
