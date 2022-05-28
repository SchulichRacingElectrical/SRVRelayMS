package main

import (
	"database-ms/app/model"
	"database-ms/config"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConfiguration struct {
	Host     string `env:"POSTGRES_HOST,required"`
	User     string `env:"POSTGRES_USER,required"`
	Password string `env:"POSTGRES_PASSWORD,required"`
	DbName   string `env:"POSTGRES_DB_NAME,required"`
	Port     string `env:"POSTGRES_PORT,required"`
	SslMode  string `env:"POSTGRES_SSLMODE,required"`
}

func main() {
	// Get the config
	conf := config.NewConfig("./env")

	// Connect to db
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		conf.Host,
		conf.User,
		conf.Password,
		conf.DbName,
		conf.Port,
		conf.SslMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(
		&model.Chart{},
		&model.ChartPreset{},
		&model.CollectionComment{},
		&model.Collection{},
		&model.Operator{},
		&model.Organization{},
		&model.RawDataPreset{},
		&model.Sensor{},
		&model.SessionComment{},
		&model.Session{},
		&model.ThingOperator{},
		&model.Thing{},
		&model.User{},
	)
}
