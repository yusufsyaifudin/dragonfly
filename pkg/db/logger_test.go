package db

import (
	"context"
	"testing"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/smartystreets/goconvey/convey"
	"go.uber.org/zap"
)

var ormDB = NewGoPGDBTest()
var z, _ = zap.NewDevelopment()

func TestGoPGDBLogger(t *testing.T) {
	convey.Convey("GoPGDBLogger", t, func() {
		convey.Convey("Debug is not active, so return immediately", func() {
			dbLogger := &GoPGDBLogger{
				Debug:     false,
				ZapLogger: z,
			}

			ctx := context.Background()

			eventQuery := &pg.QueryEvent{
				StartTime: time.Now(),
				DB:        ormDB,
				Model:     nil,
				Query:     nil,
				Params:    nil,
				Result:    nil,
				Err:       nil,
				Stash:     nil,
			}

			resCtx, err := dbLogger.BeforeQuery(ctx, eventQuery)
			convey.So(resCtx, convey.ShouldNotBeNil)
			convey.So(err, convey.ShouldBeNil)

			err = dbLogger.AfterQuery(ctx, eventQuery)
			convey.So(err, convey.ShouldBeNil)
		})

		convey.Convey("Debug is active, error occurred because no query passed in event query", func() {
			dbLogger := &GoPGDBLogger{
				Debug:     true,
				ZapLogger: z,
			}

			ctx := context.Background()

			eventQuery := &pg.QueryEvent{
				StartTime: time.Now(),
				DB:        ormDB,
				Model:     nil,
				Query:     nil, // since this query is nil, so error will be returned in this block: e.FormattedQuery()
				Params:    nil,
				Result:    nil,
				Err:       nil,
				Stash:     nil,
			}

			resCtx, err := dbLogger.BeforeQuery(ctx, eventQuery)
			convey.So(resCtx, convey.ShouldNotBeNil)
			convey.So(err, convey.ShouldNotBeNil)

			err = dbLogger.AfterQuery(ctx, eventQuery)
			convey.So(err, convey.ShouldNotBeNil)
		})

		convey.Convey("Debug is active, SELECT 1", func() {

			dbLogger := &GoPGDBLogger{
				Debug:     true,
				ZapLogger: z,
			}

			ctx := context.Background()

			eventQuery := &pg.QueryEvent{
				StartTime: time.Now(),
				DB:        ormDB,
				Model:     nil,
				Query:     "SELECT 1",
				Params:    nil,
				Result:    nil,
				Err:       nil,
				Stash:     nil,
			}

			resCtx, err := dbLogger.BeforeQuery(ctx, eventQuery)
			convey.So(resCtx, convey.ShouldNotBeNil)
			convey.So(err, convey.ShouldBeNil)

			err = dbLogger.AfterQuery(ctx, eventQuery)
			convey.So(err, convey.ShouldBeNil)
		})

		convey.Convey("Debug is active, Other query is success with args", func() {
			dbLogger := &GoPGDBLogger{
				Debug:     true,
				ZapLogger: z,
			}

			ctx := context.Background()

			eventQuery := &pg.QueryEvent{
				StartTime: time.Now(),
				DB:        ormDB,
				Model:     nil,
				Query:     "INSERT INTO users(name) VALUES (?) RETURNING *",
				Params:    []interface{}{"user name"},
				Result:    nil,
				Err:       nil,
				Stash:     nil,
			}

			resCtx, err := dbLogger.BeforeQuery(ctx, eventQuery)
			convey.So(resCtx, convey.ShouldNotBeNil)
			convey.So(err, convey.ShouldBeNil)

			err = dbLogger.AfterQuery(ctx, eventQuery)
			convey.So(err, convey.ShouldBeNil)
		})
	})
}
