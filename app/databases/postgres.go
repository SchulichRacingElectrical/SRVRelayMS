package databases

import (
	"database-ms/config"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var postgresDB *gorm.DB

func InitPostgres(config *config.Configuration) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host,
		config.User,
		config.Password,
		config.DbName,
		config.Port,
		config.SslMode,
	)
	postgresDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return postgresDB
}

func GetPostgresDB() *gorm.DB {
	return postgresDB
}
