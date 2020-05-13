package db

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/opentracing/opentracing-go"

	"github.com/go-pg/pg/v9"
)

type closer struct {
	conn *pg.DB
}

func (c closer) Close() error {
	return c.conn.Close()
}

func newCloser(conn *pg.DB) io.Closer {
	return &closer{
		conn: conn,
	}
}

// connectorGoPg using go-pg connection
func connectorGoPg(conf Conf) (*goPgDbWrapperToIFace, io.Closer, error) {
	ormPgDB := pg.Connect(&pg.Options{
		Addr:               fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		User:               conf.Username,
		Password:           conf.Password,
		Database:           conf.Database,
		ApplicationName:    conf.AppName,
		ReadTimeout:        time.Duration(conf.ReadTimeout) * time.Millisecond,
		WriteTimeout:       time.Duration(conf.WriteTimeout) * time.Millisecond,
		PoolSize:           conf.PoolSize,
		MinIdleConns:       10,
		MaxConnAge:         time.Duration(conf.MaxConnAge) * time.Millisecond,
		IdleTimeout:        time.Duration(conf.IdleTimeout) * time.Millisecond,
		IdleCheckFrequency: 1 * time.Second,
	})

	return &goPgDbWrapperToIFace{
		conf: conf,
		db:   ormPgDB,
	}, newCloser(ormPgDB), nil
}

type goPgDbWrapperToIFace struct {
	conf Conf
	db   *pg.DB
}

func (g goPgDbWrapperToIFace) ExecContext(ctx context.Context, query string, params ...interface{}) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "goPgDbWrapperToIFace.ExecContext")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	_, err := g.db.ExecContext(context.WithValue(ctx, "context", ctx), query, params...)
	return err
}

func (g goPgDbWrapperToIFace) QueryContext(ctx context.Context, model interface{}, query string, params ...interface{}) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "goPgDbWrapperToIFace.QueryContext")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	_, err := g.db.QueryContext(context.WithValue(ctx, "context", ctx), model, query, params...)
	return err
}

func (g goPgDbWrapperToIFace) Begin() (SQLTx, error) {
	tx, err := g.db.Begin()
	if err != nil {
		return nil, err
	}

	return &goPgTxWrapperToIFace{
		conf: g.conf,
		tx:   tx,
	}, nil
}

// goPgTxWrapperToIFace handling go pg transaction
type goPgTxWrapperToIFace struct {
	conf Conf
	tx   *pg.Tx
}

func (g goPgTxWrapperToIFace) ExecContext(ctx context.Context, query string, params ...interface{}) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "goPgTxWrapperToIFace.ExecContext")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	_, err := g.tx.ExecContext(context.WithValue(ctx, "context", ctx), query, params...)
	return err
}

func (g goPgTxWrapperToIFace) QueryContext(ctx context.Context, model interface{}, query string, params ...interface{}) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "goPgTxWrapperToIFace.QueryContext")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	_, err := g.tx.QueryContext(context.WithValue(ctx, "context", ctx), model, query, params...)
	return err
}

func (g goPgTxWrapperToIFace) Commit() error {
	return g.tx.Commit()
}

func (g goPgTxWrapperToIFace) Rollback() error {
	return g.tx.Rollback()
}
