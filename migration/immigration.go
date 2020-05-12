package migration

import (
	"context"
)

type Immigration interface {
	// Status get current status and list of all done migration
	Status(ctx context.Context) *CurrentStatus

	// Sync will migrate all structure into database
	Sync(ctx context.Context) error

	//Down(ctx context.Context, id int64)
}
