package tenantdb

import (
	"context"
	"fmt"

	"github.com/yusufsyaifudin/dragonfly/pkg/db"
	"github.com/yusufsyaifudin/dragonfly/tenantdb/model"

	"github.com/opentracing/opentracing-go"
)

// establishConnection establish new connection based on connection info
func establishConnection(ctx context.Context, connInfo *model.Connection) (Connection, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "establishConnection")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	if connInfo == nil {
		return nil, fmt.Errorf("connection info is nil")
	}

	conn, err := newConnection(connInfo)
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
