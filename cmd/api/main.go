package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"ysf/dragonfly/assets/migrate"
	"ysf/dragonfly/dependency"
	"ysf/dragonfly/internal/app/handler/userhandler"
	"ysf/dragonfly/pkg/db"
	"ysf/dragonfly/server"
	"ysf/dragonfly/tenantdb"
	"ysf/dragonfly/tenantdb/migration"

	"go.uber.org/zap"
)

const (
	pgHost         = "localhost"
	pgPort         = 5432
	pgUser         = "postgres"
	pgPassword     = "postgres"
	pgDB           = "dragonfly"
	pgPool         = 10
	pgIdleTimeout  = 3000  // 3 seconds
	pgMaxConnAge   = 60000 // 60 seconds
	pgReadTimeout  = 1000  // 1 seconds
	pgWriteTimeout = 2000  // 2 seconds

	appPrefix = "app"
)

func main() {
	zapLogger, err := zap.NewDevelopment(
		zap.AddCaller(),
		zap.AddCallerSkip(3),
	)

	if err != nil {
		log.Fatal(err)
		return
	}

	pgMasterConf := db.Conf{
		Disable:      false,
		Debug:        false,
		AppName:      "Dragonfly",
		Host:         pgHost,
		Port:         pgPort,
		Username:     pgUser,
		Password:     pgPassword,
		Database:     pgDB,
		PoolSize:     pgPool,
		IdleTimeout:  pgIdleTimeout,
		MaxConnAge:   pgMaxConnAge,
		ReadTimeout:  pgReadTimeout,
		WriteTimeout: pgWriteTimeout,
	}

	sqlConn, err := db.NewConnectionGoPG(db.Config{
		Master: pgMasterConf,
		Slaves: []db.Conf{},
	})

	defer func() {
		_ = sqlConn.Close()
	}()

	if err != nil {
		log.Fatal(err)
		return
	}

	// migration must per each scope,
	m := []migration.Migrate{
		new(migrate.CreateUsersTable1589341849),
	}

	// This should be initiated once in application start
	// and closed in deferred mode
	tenant, err := tenantdb.Postgres(appPrefix, sqlConn, m)
	if err != nil {
		log.Fatal("error create tenantdb connection: ", err)
		return
	}

	defer func() {
		if err := tenant.Close(); err != nil {
			fmt.Println(err.Error())
		}
	}()

	dep := dependency.NewDefault(tenant)

	// ========= Start server with graceful shutdown
	srv := server.NewServer(server.Config{
		EnableProfiling: true,
		ListenAddress:   ":2222",
		WriteTimeout:    1 * time.Second,
		ReadTimeout:     1 * time.Second,
		ZapLogger:       zapLogger,
		OpenTracing:     nil,
	})

	srv.RegisterRoutes(userhandler.Routes(dep))

	var apiErrChan = make(chan error, 1)
	go func() {
		apiErrChan <- srv.Start()
	}()

	var signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	select {
	case <-signalChan:
		_, _ = fmt.Fprintf(os.Stdout, "exiting...\n")
		srv.Shutdown()

	case err := <-apiErrChan:
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error API: %s\n", err.Error())
		}
	}
}
