package tenantdb

import (
	"context"
	"fmt"
	"ysf/dragonfly/tenantdb/model"

	"github.com/opentracing/opentracing-go"
)

func (d *pgService) getOrCreateConnectionOfTenant(ctx context.Context, tenant *model.Tenant) (connection Connection, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getOrCreateConnectionOfTenant")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	// lock read for concurrent access
	d.mutex.RLock()
	connection, exist := d.connectionById[tenant.ConnectionId]
	if exist && connection != nil {
		d.mutex.RUnlock() // unlock when return
		return
	}

	d.mutex.RUnlock() // unlock it now, no matter what

	var reader = d.conn.Reader()

	// if connection in connection is not established, establish connection to db now!
	var connInfo = &model.Connection{}
	err = reader.QueryContext(ctx, connInfo, fmt.Sprintf(sqlGetConnectionById, d.prefix), tenant.ConnectionId)
	if err != nil {
		err = fmt.Errorf("error get connection info from db: %s", err.Error())
		return
	}

	connection, err = establishConnection(ctx, connInfo)
	if err == nil && connection != nil {
		// lock and unlock for write concurrent access
		d.mutex.Lock()
		d.connectionById[connInfo.ID] = connection
		d.mutex.Unlock()
	}

	return
}
