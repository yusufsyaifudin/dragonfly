package migration

import (
	"context"
)

type Immigration interface {
	// Sync will migrate all structure into database
	Sync(ctx context.Context) error

	// Status get current status and list of all done migration
	//Status(ctx context.Context)
	//
	//Down(ctx context.Context, id int64)
}
