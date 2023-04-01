package main

import (
	"database-ms/app/model"
	"database-ms/config"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Get the config
	conf := config.NewConfig("./env")

	// Connect to default database
	connectionString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		conf.MigrationHost,
		conf.User,
		conf.Password,
		"postgres",
		conf.Port,
		conf.SslMode,
	)

	// Open connection to postgres
	postgresDB, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Check if database exists
	stmt := fmt.Sprintf("SELECT * FROM pg_database WHERE datname = '%s';", conf.DbName)
	rs := postgresDB.Raw(stmt)
	if rs.Error != nil {
		panic(rs.Error)
	}

	// Create database if not exists
	var rec = make(map[string]interface{})
	if rs.Find(rec); len(rec) == 0 {
		postgresDB.Exec("CREATE DATABASE " + conf.DbName)
	}

	// Connect to db
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		conf.MigrationHost,
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

	// Full migration
	db.AutoMigrate(
		&model.Blacklist{},
		&model.Chart{},
		&model.ChartPreset{},
		&model.Collection{},
		&model.Datum{},
		&model.Operator{},
		&model.Organization{},
		&model.RawDataPreset{},
		&model.Sensor{},
		&model.Session{},
		&model.ThingOperator{},
		&model.Thing{},
		&model.User{},
		&model.Comment{},
		&model.SessionCollection{},
		&model.RawDataPresetSensor{},
		&model.ChartSensor{},
	)

	println("Finished migration.")
}
