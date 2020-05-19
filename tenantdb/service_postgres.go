package tenantdb

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"
	"ysf/dragonfly/pkg/db"
	"ysf/dragonfly/tenantdb/migration"
	"ysf/dragonfly/tenantdb/model"

	"github.com/pkg/errors"

	"github.com/opentracing/opentracing-go"
)

var (
	ErrNoConnectionList = fmt.Errorf("connections list is empty, could not select one")
)

// pgService implements Service using Postgres as main database to save all tenant data and it's connection info.
type pgService struct {
	mutex          sync.RWMutex
	prefix         string
	migrates       []migration.Migrate
	conn           db.SQL
	connectionById map[string]Connection // connection id -> connection
}

func (d pgService) CreateTenant(ctx context.Context, tenantId, tenantName string) (connection Connection, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateTenant")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	tenantId, err = SanitizeTenantId(ctx, tenantId)
	if err != nil {
		return
	}

	// TODO: do the least usage of db connection list

	var writer = d.conn.Writer()

	var connections = model.Connections{}
	err = writer.QueryContext(ctx, &connections, fmt.Sprintf(sqlGetConnections, d.prefix))
	if err != nil {
		return
	}

	if len(connections) <= 0 {
		err = ErrNoConnectionList
		return
	}

	var connInfo = &model.Connection{}
	for _, conn := range connections {
		connInfo = conn
	}
	var tenant = &model.Tenant{}
	_ = writer.QueryContext(ctx, tenant, fmt.Sprintf(sqlGetTenantByID, d.prefix), tenantId)
	if tenant.ID == "" {
		err = writer.QueryContext(ctx, tenant, fmt.Sprintf(sqlCreateTenant, d.prefix), tenantId, tenantName, connInfo.ID)
	}

	if err != nil {
		return
	}

	connection, err = d.getOrCreateConnectionOfTenant(ctx, tenant)
	return
}

func (d pgService) GetTenant(ctx context.Context, tenantId string) (connection Connection, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetTenant")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	tenantId, err = SanitizeTenantId(ctx, tenantId)
	if err != nil {
		return
	}

	var reader = d.conn.Reader()

	var tenant = &model.Tenant{}
	err = reader.QueryContext(ctx, tenant, fmt.Sprintf(sqlGetTenantByID, d.prefix), tenantId)
	if err != nil {
		return
	}

	if tenant.ID == "" {
		err = fmt.Errorf("tenant id %s is not found", tenantId)
		return
	}

	connection, err = d.getOrCreateConnectionOfTenant(ctx, tenant)
	return
}

func (d pgService) GetTenants(ctx context.Context) (tenants model.Tenants, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetTenants")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	tenants = model.Tenants{}
	err = d.conn.Reader().QueryContext(ctx, &tenants, fmt.Sprintf(sqlGetTenants, d.prefix))
	return
}

func (d pgService) GetTenantImmigration(ctx context.Context, tenantId string) (immigration migration.Immigration, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetTenantImmigration")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	tenantId, err = SanitizeTenantId(ctx, tenantId)
	if err != nil {
		return
	}

	var reader = d.conn.Reader()

	var tenant = &model.Tenant{}
	err = reader.QueryContext(ctx, tenant, fmt.Sprintf(sqlGetTenantByID, d.prefix), tenantId)
	if err != nil {
		return
	}

	if tenant.ID == "" {
		err = fmt.Errorf("tenant id %s is not found", tenantId)
		return
	}

	connection, err := d.getOrCreateConnectionOfTenant(ctx, tenant)
	if err != nil {
		return
	}

	if connection.SQL() == nil {
		return migration.NewNoopImmigration(), fmt.Errorf("connection to sql is nil")
	}

	return migration.NewImmigrationPostgres(connection.SQL(), fmt.Sprintf("%s_%s", d.prefix, tenantId), d.migrates)
}

func (d pgService) getOrCreateConnectionOfTenant(ctx context.Context, tenant *model.Tenant) (connection Connection, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getOrCreateConnectionOfTenant")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	d.mutex.Lock()
	defer d.mutex.Unlock()
	connection, exist := d.connectionById[tenant.ConnectionId]
	if exist && connection != nil {
		return
	}

	var reader = d.conn.Reader()

	// if connection in connection is not established, establish connection to db now!
	var connInfo = &model.Connection{}
	err = reader.QueryContext(ctx, connInfo, fmt.Sprintf(sqlGetConnectionById, d.prefix), tenant.ConnectionId)
	if err != nil {
		return
	}

	connection, err = establishConnection(ctx, tenant, connInfo)
	if err == nil && connection != nil {
		d.connectionById[connInfo.ID] = connection
	}

	return
}

func (d pgService) Close() error {
	var err error
	for _, c := range d.connectionById {
		if c.SQL() == nil {
			continue
		}

		sqlErr := c.SQL().Close()
		if sqlErr != nil {
			err = errors.Wrapf(err, sqlErr.Error())
		}
	}

	return err
}

func Postgres(prefix string, conn db.SQL, migrates []migration.Migrate) (service Service, err error) {
	ctx := context.Background()

	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		err = fmt.Errorf("prefix is empty")
		return
	}

	var prefixBuf = bytes.Buffer{}
	for _, char := range prefix {
		_, isAlphabet := lowercaseAlphabetChars[char]

		if isAlphabet {
			prefixBuf.WriteRune(char)
		}
	}

	if prefixBuf.String() != prefix {
		err = fmt.Errorf("prefix must be lowercase alphabet only")
		return
	}

	prefix = prefixBuf.String()

	service = &pgService{}

	writer := conn.Writer()
	tx, err := writer.BeginTx()
	if err != nil {
		return
	}

	// do all migration at once
	err = tx.ExecContext(ctx, sqlCreateUuidExt)
	if err != nil {
		_ = tx.Rollback()
		return
	}

	err = tx.ExecContext(ctx, fmt.Sprintf(sqlCreateConnectionTable, prefix))
	if err != nil {
		_ = tx.Rollback()
		return
	}

	err = tx.ExecContext(ctx, fmt.Sprintf(
		sqlCreateTenantsTable,
		prefix,
		prefix,
		prefix,
		prefix,
		prefix,
	))
	if err != nil {
		_ = tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	return &pgService{
		mutex:          sync.RWMutex{},
		prefix:         prefix,
		migrates:       migrates,
		conn:           conn,
		connectionById: make(map[string]Connection),
	}, nil
}
