# Overview
Hecate, Greek Goddess guarding of the portal between realms, is the one stop shop for planning your next trip

## Development

1. Download latest release of Go
2. Start postgres container with `docker-compose up -d`
3. `go run *.go` builds and runs the final binary

## Database backups

First install `psql` and then run:

```
pg_dump -h localhost -U admin hecate -Fc > db-backups/"$(date)".sql
```

To restore:

```
pg_restore -h localhost -d hecate -U admin -C "name"
```

To sync to S3:

```
aws s3 sync db-backups s3://hecate-backups
```

