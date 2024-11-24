# Overview
Hecate, Greek Goddess guarding of the portal between realms, is the one stop shop for planning your next trip

## Development

1. Download latest release of Go
2. Start postgres container with `docker-compose up -d`
3. `go run *.go` builds and runs the final binary

## Database backups (Automated)

The docker compose uses another sidecar service to run backups of the postgres database. The backups are stored in db-backups/monthly|weekly|daily|last

### To restore from one of these backups

```
cat hecate-20241124-180000.sql.gz  | gunzip | docker exec -i hecate-db-1 psql -U admin -d hecate
```

## Database backups (Manual)

First install `psql` and then run:

```
pg_dump -h localhost -U admin hecate -Fc > db-backups/"$(date)".sql
```

To restore:

```
cat your_dump.sql | docker exec -i your-db-container psql -U admin
```

Or if running locally:

```
pg_restore -h localhost -d hecate -U admin -C "name"
```

To sync to S3:

```
aws s3 sync db-backups s3://hecate-backups
```

