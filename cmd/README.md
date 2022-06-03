## PostgreSQL DB with GORM

## Prerequisite

- Postgresql version 10
- pgAdmin 4

## DB Setup

1. Connect to the local database using pgAdmin
2. Create a database named `sr-velocity`

**Steps**

1. Duplicate the .env file file and place it in the `/cmd` folder
2. Run: `go run cmd/migrate.go`
3. The database should be updated

Note that columns will not be deleted if they are removed from the schema.
