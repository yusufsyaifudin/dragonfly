package tenantdb

import (
	"ysf/dragonfly/pkg/db"
	"ysf/dragonfly/tenantdb/model"
)

// establishConnection establish new connection based on connection info
func establishConnection(tenant *model.Tenant, connInfo *model.Connection) (Connection, error) {
	conn := NewConnection(tenant, connInfo)

	pgMaster, err := db.NewConnectionGoPG(db.Config{
		Master: db.Conf{
			Disable:      false,
			Debug:        connInfo.PostgresMasterDebug,
			AppName:      connInfo.ID,
			Host:         connInfo.PostgresMasterHost,
			Port:         connInfo.PostgresMasterPort,
			Username:     connInfo.PostgresMasterUsername,
			Password:     connInfo.PostgresMasterPassword,
			Database:     connInfo.PostgresMasterDatabase,
			PoolSize:     connInfo.PostgresMasterPoolSize,
			IdleTimeout:  connInfo.PostgresMasterIdleTimeout,
			MaxConnAge:   connInfo.PostgresMasterMaxConnAge,
			ReadTimeout:  connInfo.PostgresMasterReadTimeout,
			WriteTimeout: connInfo.PostgresMasterWriteTimeout,
		},
		Slaves: nil,
	})

	if err != nil {
		return nil, err
	}

	// TODO: connect to redis

	conn.setSQL(pgMaster)

	return conn, nil
}
