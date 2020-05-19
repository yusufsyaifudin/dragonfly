# Dragonfly

> Multi Tenant Skeleton App (WIP)

This repo gives you a basic structure how design multi-tenant database app.
You can write a migration once, then using argument `tenantId` in each Up and Down migration function, you can create the
schema in Postgre based on tenant.id in database.

Best practice:
* Please separate the connection info and tenants data in separate database (work in progress to make it work)
* Use tenant.id as schema name in postgres so it will share the connection pool (that's why I just save database name in connection info, not schema name)
* Validate and sanitize each tenant.id before saving into database, so it conform with table name rule in every database (don't use number as a prefix)

## Variables

Set this (also in app_connections table too):

```
    pgHost         = "localhost"
	pgPort         = 5432
	pgUser         = "postgres"
	pgPassword     = "postgres"
	pgDB           = "dragonfly"
	pgPool         = 30
	pgReadTimeout  = 1000  // 1 seconds
	pgWriteTimeout = 2000  // 2 seconds

	appPrefix = "app"
```

```
INSERT INTO public.app_connections (postgres_master_debug,postgres_master_host,postgres_master_port,postgres_master_username,postgres_master_password,postgres_master_database,postgres_master_pool_size,postgres_master_read_timeout,postgres_master_write_timeout) VALUES 
(true,'localhost',5432,'postgres','postgres','dragonfly',30,1000,2000);
```