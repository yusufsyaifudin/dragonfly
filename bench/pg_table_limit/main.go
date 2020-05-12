package main

import (
	"fmt"
	"os"

	"github.com/go-pg/pg/v9"
)

const (
	username = "postgres"
	password = "postgres"
	database = "test_db"
	address  = "localhost"
	port     = "5433"

	numTable  = 1_000_000
	numWorker = 80 // max pool size default on postgres is 100
)

func main() {
	db := pg.Connect(&pg.Options{
		User:     username,
		Password: password,
		Database: database,
		Addr:     fmt.Sprintf("%s:%s", address, port),
	})

	defer func() {
		if err := db.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error closing db %s\n", err.Error())
		}
	}()

	// do this before running the queries:
	// CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	var createTableQuery = func(i int) string {
		sql := `
			CREATE TABLE IF NOT EXISTS event_%d (
				sequence_num BIGINT generated always AS IDENTITY,
				tenant_id UUID NOT NULL,
				stream_id UUID NOT NULL,
				version int NOT NULL DEFAULT 1,
				type text NOT NULL DEFAULT '',
				meta jsonb NOT NULL DEFAULT '{}',
				data jsonb NOT NULL DEFAULT '{}',
				log_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
				CONSTRAINT pk_event_%d_sequence_num PRIMARY KEY (sequence_num, tenant_id),
				CONSTRAINT uk_event_%d_stream_id_version UNIQUE (tenant_id, stream_id, version)
			);
		`
		return fmt.Sprintf(sql, i, i, i)
	}

	var worker = func(id int, jobs <-chan int, errChan chan<- error) {
		for j := range jobs {
			_, _ = fmt.Fprintf(os.Stdout, "doing worker %d job %d\n", id, j)

			_, err := db.Exec(createTableQuery(j))
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "error create table %d %s\n", j, err.Error())
				errChan <- err
			}
		}
	}

	jobs := make(chan int, numTable)
	results := make(chan error, numTable)

	for w := 1; w <= numWorker; w++ {
		go worker(w, jobs, results)
	}

	for j := 346418; j <= numTable; j++ {
		jobs <- j
	}

	close(jobs)

	for a := 346418; a <= numTable; a++ {
		<-results
	}

}
