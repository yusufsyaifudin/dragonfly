package tenantdb

import (
	"github.com/yusufsyaifudin/dragonfly/pkg/db"
	"github.com/yusufsyaifudin/dragonfly/tenantdb/model"
)

// Connection real connection to specific tenant
type Connection interface {
	ConnectionInfo() *model.Connection
	SQL() db.SQL
	Redis()
}
