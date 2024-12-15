# Overview
Hecate, Greek Goddess guarding of the portal between realms, is the one stop shop for planning your next trip

## Development

1. Download latest release of Go
2. `docker-compose up -d` builds and runs the application
3. `go run *.go` builds and runs the final binary locally

## Database

The project now uses SQLite for local data storage. The database file is automatically created in the `/app/data` directory when the application starts.

### Migrating from PostgreSQL

If you have an existing PostgreSQL database dump and want to migrate to SQLite:

1. Ensure you have `pg_restore` and `sqlite3` installed
2. Run the migration script:
   ```bash
   ./scripts/migrate_postgres_to_sqlite.sh path/to/your/postgres_dump.sql
   ```

#### Migration Script Details

The migration script does the following:
- Converts PostgreSQL custom dump to plain SQL
- Modifies SQL to be SQLite-compatible
  - Replaces SERIAL with INTEGER PRIMARY KEY AUTOINCREMENT
  - Converts JSONB to TEXT
  - Adjusts sequence and data type references
- Creates a new SQLite database file
- Imports the converted data

### Database Persistence

- The SQLite database is persisted using a Docker volume `sqlite_data`
- The database file is located at `/app/data/hecate.db` inside the container
- When running locally, the database will be created in a `./data` directory

### Backup and Migration

Since SQLite is a file-based database, you can backup the database by simply copying the `.db` file:

```bash
# Inside the container
docker exec hecate cp /app/data/hecate.db /app/data/hecate_backup.db

# Or locally
cp ./data/hecate.db ./data/hecate_backup.db
```

To migrate the database to a new system, just copy the `.db` file to the new location.

## Development Workflow

1. Start the application:
   ```bash
   docker-compose up -d
   ```

2. Stop the application:
   ```bash
   docker-compose down
   ```

3. Rebuild the application:
   ```bash
   docker-compose up -d --build
   ```

## Local Development

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Run the application:
   ```bash
   go run *.go
