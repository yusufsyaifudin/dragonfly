package migration

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"ysf/dragonfly/pkg/db"

	"github.com/opentracing/opentracing-go"
)

type postgres struct {
	master    db.SQLWriter
	tenantId  string
	migration []Migrate
}

func (p postgres) Status(ctx context.Context) *CurrentStatus {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Status")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	var tenantID = p.tenantId

	// order ascending by sequence number
	var appliedMigrations = make([]*ModelMigration, 0)
	err := p.master.QueryContext(ctx, &appliedMigrations, fmt.Sprintf(sqlGetMigrationData, tenantID))
	if err != nil {
		return &CurrentStatus{}
	}

	lastApplied := appliedMigrations[len(appliedMigrations)-1]
	if lastApplied == nil {
		lastApplied = &ModelMigration{}
	}

	return &CurrentStatus{
		AppliedMigrations:  appliedMigrations,
		LastID:             lastApplied.ID,
		LastSequenceNumber: lastApplied.SequenceNumber,
		LastAppliedAt:      lastApplied.AppliedAt,
	}
}

// Sync do migration to specific postgres schema.
// Time complexity: O(3 + 2N) where
// 3 is minimum query to create schema, create migration table, and get applied migration
// 2 is fixed number when we sync up the migration. For example, when we do create users table,
// it will first create the table in the selected schema and then insert the record to applied migrations data.
// N is variable, the number of migration file to be sync.
func (p postgres) Sync(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Sync")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	var tenantID = p.tenantId

	// order ascending by sequence number
	var migrations = p.migration
	sort.SliceStable(migrations, func(i, j int) bool {
		return migrations[i].SequenceNumber(ctx) < migrations[j].SequenceNumber(ctx)
	})

	var queries = bytes.Buffer{}
	defer queries.Reset()

	queries.WriteString(fmt.Sprintf(sqlCreatePostgresSchema, tenantID))
	queries.WriteString(fmt.Sprintf(sqlCreateMigrationTable, tenantID, tenantID, tenantID))
	queries.WriteString(fmt.Sprintf(sqlGetMigrationData, tenantID))

	// since if it frequently called, and doesn't need rollback, then use preparedCreateTable statement
	var appliedMigrations = make([]*ModelMigration, 0)
	err := p.master.QueryContext(ctx, &appliedMigrations, queries.String())
	if err != nil {
		return fmt.Errorf("error get migration data: %s", err.Error())
	}

	if len(appliedMigrations) > len(migrations) {
		return fmt.Errorf("some applied migration in database is not found in list of migration")
	}

	var candidateMigration = make(map[string]bool, 0)
	for _, m := range migrations {
		candidateMigration[m.ID(ctx)] = true
	}

	var appliedMigrationById = make(map[string]*ModelMigration)
	for _, m := range appliedMigrations {
		if _, exist := candidateMigration[m.ID]; !exist {
			return fmt.Errorf("migration %s is not registred in migration list", m.ID)
		}

		appliedMigrationById[m.ID] = m
	}

	type up struct {
		id          string
		sequenceNum int
		query       string
	}

	var candidateUp = make([]up, 0)
	for _, m := range migrations {
		id := m.ID(ctx)
		seqNum := m.SequenceNumber(ctx)

		if v, exist := appliedMigrationById[id]; exist {
			// detect anomalies, when sequence number in input data is not same as order
			// in applied migration database
			if v.SequenceNumber != seqNum {
				return fmt.Errorf(
					"%s is registered with sequence %d but applied with sequence %d",
					v.ID,
					seqNum,
					v.SequenceNumber,
				)
			}

			continue
		}

		sql, err := m.Up(ctx, tenantID)
		if err != nil {
			return err
		}

		candidateUp = append(candidateUp, up{
			id:          id,
			sequenceNum: seqNum,
			query:       sql,
		})
	}

	// if no migration to be synced, then early return without error
	if len(candidateUp) <= 0 {
		return nil
	}

	tx := p.master

	for _, m := range candidateUp {
		err = tx.ExecContext(ctx, m.query)
		if err != nil {
			return fmt.Errorf("error execution migration id %s: %s", m.id, err.Error())
		}

		err = tx.ExecContext(ctx, fmt.Sprintf(sqlInsertMigrationData, tenantID), m.id, m.sequenceNum)
		if err != nil {
			return fmt.Errorf("error recording %s: %s", m.id, err.Error())
		}
	}

	return nil
}

func NewImmigrationPostgres(conn db.SQL, tenantId string, migration []Migrate) (Immigration, error) {
	return &postgres{
		master:    conn.Writer(),
		tenantId:  tenantId,
		migration: migration,
	}, nil
}
