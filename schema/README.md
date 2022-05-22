## PostgreSQL DB with GORM

## Prerequisite

- Postgresql version 10
- pgAdmin 4

## DB Setup

1. Connect to the local database using pgAdmin
2. Create a database named `sr-velocity`
3. Open the query tool and copy paste the query in `schema/rdb_init.sql`
4. Run the query

## Making Schema Changes

The end goal would be do use GORM's AutoMigration feature, but to get things working quickly, we are using GORM Gen to generate the models

**Steps**

1. Make the changes to the database schema(use pgAdmin to make this process easier)
2. Verify that that changes are applied
3. Run the generate script `cmd/generate.go` by running the command `go run generate.go`. You might need to change the password in the dsn string
4. There should be a the generated model files in `app/model` directory
