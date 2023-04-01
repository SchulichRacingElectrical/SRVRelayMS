package databases

import (
	"database-ms/config"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var postgresDB *gorm.DB

func InitPostgres(config *config.Configuration) *gorm.DB {
	connectionString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host,
		config.User,
		config.Password,
		"postgres",
		config.Port,
		config.SslMode,
	)

	// Open connection to postgres
	postgresDB, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Check if database exists
	stmt := fmt.Sprintf("SELECT * FROM pg_database WHERE datname = '%s';", config.DbName)
	rs := postgresDB.Raw(stmt)
	if rs.Error != nil {
		panic(rs.Error)
	}

	// Create database if not exists
	var rec = make(map[string]interface{})
	if rs.Find(rec); len(rec) == 0 {
		postgresDB.Exec("CREATE DATABASE " + config.DbName)
	}

	// Get new connection to database
	connectionString = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host,
		config.User,
		config.Password,
		config.DbName,
		config.Port,
		config.SslMode,
	)
	postgresDB, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return postgresDB
}

func GetPostgresDB() *gorm.DB {
	return postgresDB
}
