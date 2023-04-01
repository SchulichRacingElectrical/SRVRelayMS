## PostgreSQL DB with GORM

**Steps**

1. Run `docker compose up` to start the database
2. Run: `go run cmd/migrate.go`
3. The database should be updated

Note that columns will not be deleted if they are removed from the schema.
