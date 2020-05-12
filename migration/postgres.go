package migration

import (
	"context"
	"ysf/dragonfly/pkg/db"
)

type postgres struct {
	master    db.SQLWriter
	tenantId  string
	migration []Migrate
}

func (p postgres) Sync(ctx context.Context) error {
	type up struct {
		id    string
		query string
	}

	var candidateUp = make([]up, 0)
	for _, m := range p.migration {
		sql, err := m.Up(ctx, p.tenantId)
		if err != nil {
			return err
		}

		candidateUp = append(candidateUp, up{
			id:    m.ID(ctx),
			query: sql,
		})
	}

	tx, err := p.master.Begin()
	if err != nil {
		return err
	}

	// TODO create migration table

	for _, m := range candidateUp {
		err = tx.ExecContext(ctx, m.query)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		// TODO: insert into db to handle query
		// INSERT INTO tenantID_migrations (id, applied_at) VALUES (m.id, now())

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
