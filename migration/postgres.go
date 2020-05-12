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

	var appliedMigrations = make([]*ModelMigration, 0)
	err := p.master.QueryContext(ctx, &appliedMigrations, fmt.Sprintf(sqlGetMigrationData, tenantID, tenantID))
	if err != nil {
		return &CurrentStatus{}
	}

	// order ascending by sequence number
	sort.SliceStable(appliedMigrations, func(i, j int) bool {
		return appliedMigrations[i].SequenceNumber < appliedMigrations[j].SequenceNumber
	})

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

	var appliedMigrations = make([]*ModelMigration, 0)
	err := p.master.QueryContext(ctx, &appliedMigrations, fmt.Sprintf(sqlGetMigrationData, tenantID, tenantID))
	if err != nil {
		return err
	}

	var appliedMigrationById = make(map[string]bool)
	for _, m := range appliedMigrations {
		appliedMigrationById[m.ID] = true
	}

	type up struct {
		id          string
		sequenceNum int
		query       string
	}

	// TODO: detect anomalies, when order in input data is not same as order in applied migration database
	var candidateUp = make([]up, 0)
	for _, m := range p.migration {
		if ok, exist := appliedMigrationById[m.ID(ctx)]; ok && exist {
			continue
		}

		sql, err := m.Up(ctx, tenantID)
		if err != nil {
			return err
		}

		candidateUp = append(candidateUp, up{
			id:          m.ID(ctx),
			sequenceNum: m.SequenceNumber(ctx),
			query:       sql,
		})
	}

	// if no migration to be synced, then early return without error
	if len(candidateUp) <= 0 {
		return nil
	}

	// order ascending by sequence number
	sort.SliceStable(candidateUp, func(i, j int) bool {
		return candidateUp[i].sequenceNum < candidateUp[j].sequenceNum
	})

	tx, err := p.master.Begin()
	if err != nil {
		return err
	}

	err = tx.ExecContext(ctx, fmt.Sprintf(sqlCreateMigrationTable, tenantID, tenantID))
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
