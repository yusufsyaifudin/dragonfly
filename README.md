# Dragonfly

> Multi Tenant Skeleton App (WIP)

This repo gives you a basic structure how design multi-tenant database app.
You can write a migration once, then using argument `tenantId` in each Up and Down migration function, you can create the
schema in Postgre based on tenant.id in database.

Best practice:
* Please separate the connection info and tenants data in separate database (work in progress to make it work)
* Use tenant.id as schema name in postgres so it will share the connection pool (that's why I just save database name in connection info, not schema name)
* Validate and sanitize each tenant.id before saving into database, so it conform with table name rule in every database (don't use number as a prefix)