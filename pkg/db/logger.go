package db

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/opentracing/opentracing-go"
	field "github.com/opentracing/opentracing-go/log"
	"go.uber.org/zap"
)

// GoPGDBLogger Logger for go db
type GoPGDBLogger struct {
	Debug     bool
	ZapLogger *zap.Logger
}

func (l *GoPGDBLogger) BeforeQuery(ctx context.Context, e *pg.QueryEvent) (context.Context, error) {
	if ctx == nil {
		ctx = e.DB.Context()
	}

	if ctx == nil {
		ctx = context.Background()
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "BeforeQuery")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	if !l.Debug {
		return ctx, nil
	}

	query, err := e.FormattedQuery()
	if err != nil {
		query = fmt.Sprintf("query cannot builded due to error: %s", err.Error())
		return ctx, err
	}

	// replace duplicate space
	space := regexp.MustCompile(`\s+`)
	query = strings.TrimSpace(space.ReplaceAllString(query, " "))

	span.LogFields(
		field.String("query", query),
	)

	return ctx, nil
}

func (l *GoPGDBLogger) AfterQuery(ctx context.Context, e *pg.QueryEvent) error {
	if ctx == nil {
		ctx = e.DB.Context()
	}

	if ctx == nil {
		ctx = context.Background()
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "AfterQuery")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	if !l.Debug {
		return nil
	}

	query, err := e.FormattedQuery()
	if err != nil {
		query = fmt.Sprintf("query cannot builded due to error: %s", err.Error())
		return err
	}

	st := e.StartTime
	elapsedTime := float64(time.Since(st).Nanoseconds()) / float64(1000000) // in ms

	if l.Debug {
		l.ZapLogger.Info(
			"DBQuery",
			zap.Float64("elapsed_time", elapsedTime),
			zap.String("query", query),
		)

	}

	span.LogFields(
		field.String("elapsedTime", fmt.Sprintf("%f ms", elapsedTime)),
	)

	// replace duplicate space
	space := regexp.MustCompile(`\s+`)
	query = strings.TrimSpace(space.ReplaceAllString(query, " "))

	return nil
}
