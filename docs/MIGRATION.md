# Database Migration Guide

This repository does not include a migration framework by default. Use the following pattern for Postgres migrations.

## Recommended approach

1. Create a `migrations/` directory at the repository root.
2. Add ordered SQL migration files, for example:
   - `0001_create_schema.sql`
   - `0002_add_maintenance_tables.sql`
3. Track migrations in source control.

## Using `psql`

Apply migrations manually:

```powershell
$env:PGUSER='<your-user>'
$env:PGPASSWORD='<your-password>'
$env:PGHOST='<your-host>'
$env:PGPORT='5432'
$env:PGDATABASE='<your-db>'

psql "host=$env:PGHOST port=$env:PGPORT user=$env:PGUSER password=$env:PGPASSWORD dbname=$env:PGDATABASE sslmode=disable" -f migrations/0001_create_schema.sql
```

## Example migration steps

### 1. Create schema

```sql
CREATE TABLE IF NOT EXISTS vendors (
  id SERIAL PRIMARY KEY,
  uid UUID NOT NULL UNIQUE,
  name TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW()
);
```

### 2. Add maintenance schedule tables

```sql
CREATE TABLE IF NOT EXISTS maintenance_schedules (
  id SERIAL PRIMARY KEY,
  uid UUID NOT NULL UNIQUE,
  maintenance_template_id INTEGER NOT NULL,
  machine_id INTEGER NOT NULL,
  start_timestamp TIMESTAMPTZ,
  end_timestamp TIMESTAMPTZ
);
```

## Recommended migration tools

If you want automated migration tracking, add one of the following:

- `golang-migrate/migrate`
- `pressly/goose`
- `rubenv/sql-migrate`

Example using `migrate`:

```powershell
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate -source file://migrations -database "postgres://$env:PGUSER:$env:PGPASSWORD@$env:PGHOST:$env:PGPORT/$env:PGDATABASE?sslmode=disable" up
```

## Environment promotion for migrations

- DEV: run migrations first in a disposable dev database.
- UAT: run migrations against a copy of the production schema before promotion.
- PROD: apply migrations during a controlled maintenance window.

## Backup and rollback

- Take a database backup before applying production migrations.
- Use `pg_dump` to export the current schema and data.

```powershell
pg_dump --format=custom --file=backup_$(Get-Date -Format yyyyMMddHHmmss).dump
```

- Rollback by restoring from the backup if a migration fails.
