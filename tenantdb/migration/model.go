package migration

import (
	"time"
)

// ModelMigration represent each migration table name
type ModelMigration struct {
	ID             string
	SequenceNumber int
	AppliedAt      time.Time
}

type CurrentStatus struct {
	AppliedMigrations  []*ModelMigration
	LastID             string
	LastSequenceNumber int
	LastAppliedAt      time.Time
}
