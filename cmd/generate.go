package main

import (
	"fmt"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	// Connect to db
	dsn := "host=localhost user=postgres password=foobar123 dbname=sr-velocity port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// get file location
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println(path) // for example /home/user

	g := gen.NewGenerator(gen.Config{
		OutPath: strings.Replace(path, "cmd", "app/models", -1),
	})

	g.UseDB(db)

	g.GenerateAllTable()

	g.Execute()
}
