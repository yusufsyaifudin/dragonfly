package db

import (
	"context"
)

type noopTx struct{}

func (n noopTx) ExecContext(_ context.Context, _ string, _ ...interface{}) error { return nil }

func (n noopTx) QueryContext(_ context.Context, _ interface{}, _ string, _ ...interface{}) error {
	return nil
}

func (n noopTx) Commit() error { return nil }

func (n noopTx) Rollback() error { return nil }

type noopWriter struct{}

func (n noopWriter) ExecContext(_ context.Context, _ string, _ ...interface{}) error { return nil }

func (n noopWriter) QueryContext(_ context.Context, _ interface{}, _ string, _ ...interface{}) error {
	return nil
}

func (n noopWriter) Begin() (SQLTx, error) {
	return &noopTx{}, nil
}

func (n noopSql) Writer() SQLWriter {
	return &noopWriter{}
}

func (n noopSql) Reader() SQLReader {
	return &noopWriter{}
}

type noopSql struct{}

func (n noopSql) Close() error { return nil }

func NewNoop() SQL {
	return &noopSql{}
}
