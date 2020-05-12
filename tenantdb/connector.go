package tenantdb

import (
	"context"
	"ysf/dragonfly/pkg/db"
	"ysf/dragonfly/tenantdb/model"

	"github.com/opentracing/opentracing-go"
)

// establishConnection establish new connection based on connection info
func establishConnection(ctx context.Context, tenant *model.Tenant, connInfo *model.Connection) (Connection, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "establishConnection")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	conn, err := newConnection(tenant, connInfo)
	if err != nil {
		return nil, err
	}

	pgMasterConf := db.Conf{
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
	}

	pgMaster, err := db.NewConnectionGoPG(db.Config{
		Master: pgMasterConf,
		Slaves: []db.Conf{
			pgMasterConf,
		},
	})

	if err != nil {
		return nil, err
	}

	// TODO: connect to redis

	conn.setSQL(pgMaster)

	return conn, nil
}
