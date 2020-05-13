package db

import (
	"context"
)

// SQL can be use writer or read replica only.
type SQL interface {
	Writer() SQLWriter
	Reader() SQLReader
	Close() error
}

// SQLWriter is a common interface for pg.DB and pg.Tx types.
// this copies from go-pg DB interface
type SQLWriter interface {
	ExecContext(c context.Context, query string, params ...interface{}) error
	QueryContext(c context.Context, model interface{}, query string, params ...interface{}) error

	Begin() (SQLTx, error)
}

// SQLTx is a interface to handle go-pg orm transaction
type SQLTx interface {
	ExecContext(c context.Context, query string, params ...interface{}) error
	QueryContext(c context.Context, model interface{}, query string, params ...interface{}) error
	Commit() error
	Rollback() error
}

// SQLReader use slave instance
type SQLReader interface {
	ExecContext(c context.Context, query string, params ...interface{}) error
	QueryContext(c context.Context, model interface{}, query string, params ...interface{}) error
}
