package migration

import "context"

type noopImmigration struct{}

func (n noopImmigration) Sync(_ context.Context) error { return nil }

func (n noopImmigration) Close() error { return nil }

func NewNoopImmigration() Immigration {
	return &noopImmigration{}
}
