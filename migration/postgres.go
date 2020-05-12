package migration

import (
	"context"
	"fmt"
	"sort"
	"ysf/dragonfly/pkg/db"
)

type postgres struct {
	master    db.SQLWriter
	tenantId  string
	migration []Migrate
}

func (p postgres) Status(ctx context.Context) *CurrentStatus {
	var tenantID = p.tenantId

	// order ascending by sequence number
	var appliedMigrations = make([]*ModelMigration, 0)
	err := p.master.QueryContext(ctx, &appliedMigrations, fmt.Sprintf(sqlGetMigrationData, tenantID, tenantID))
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

func (p postgres) Sync(ctx context.Context) error {
	var tenantID = p.tenantId

	// order ascending by sequence number
	var migrations = p.migration
	sort.SliceStable(migrations, func(i, j int) bool {
		return migrations[i].SequenceNumber(ctx) < migrations[j].SequenceNumber(ctx)
	})

	var appliedMigrations = make([]*ModelMigration, 0)
	err := p.master.QueryContext(ctx, &appliedMigrations, fmt.Sprintf(sqlGetMigrationData, tenantID, tenantID))
	if err != nil {
		return err
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
			// detect anomalies, when sequence number in input data is not same as order in applied migration database
			if v.SequenceNumber != seqNum {
				return fmt.Errorf("%s is registered with sequence %d but applied with sequence %d",
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

	tx, err := p.master.Begin()
	if err != nil {
		return err
	}

	err = tx.ExecContext(ctx, fmt.Sprintf(sqlCreateMigrationTable, tenantID, tenantID, tenantID, tenantID, tenantID))
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	for _, m := range candidateUp {
		err = tx.ExecContext(ctx, m.query)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		err = tx.ExecContext(ctx, fmt.Sprintf(sqlInsertMigrationData, tenantID, tenantID), m.id, m.sequenceNum)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func NewPostgres(conn db.SQL, tenantId string, migration []Migrate) Immigration {
	return &postgres{
		master:    conn.Writer(),
		tenantId:  tenantId,
		migration: migration,
	}
}
