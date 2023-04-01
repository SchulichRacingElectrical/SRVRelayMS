package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Configuration struct {
	Host          string `env:"POSTGRES_HOST,required"`
	MigrationHost string `env:"POSTGRES_MIGRATION_HOST,required"`
	User          string `env:"POSTGRES_USER,required"`
	Password      string `env:"POSTGRES_PASSWORD,required"`
	DbName        string `env:"POSTGRES_DB_NAME,required"`
	Port          string `env:"POSTGRES_PORT,required"`
	SslMode       string `env:"POSTGRES_SSLMODE,required"`
	AdminKey      string `env:"ADMIN_API_KEY,required"`
	AccessSecret  string `env:"ACCESS_SECRET,required"`
	Address       string `env:"ADDRESS" envDefault:":8000"`
	RedisUrl      string `env:"REDIS_URL,required"`
	RedisPort     string `env:"REDIS_PORT,required"`
	RedisUsername string `env:"REDIS_USERNAME,required"`
	RedisPassword string `env:"REDIS_PASSWORD,required"`
	FilePath      string `env:"FILE_PATH,required"`
}

// NewConfig will read the config data from given .env file
func NewConfig(files ...string) *Configuration {
	path, _ := os.Getwd()
	fullpath := filepath.Join(path, ".env")
	fmt.Println(fullpath)
	err := godotenv.Load(fullpath) // Loading config from env file

	if err != nil {
		log.Printf("No .env file could be found %q\n", files)
	}

	cfg := Configuration{}

	// Parse env to configuration
	err = env.Parse(&cfg)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	return &cfg
}
