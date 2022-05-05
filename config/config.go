package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

// Configuration contains static info required to run the apps
type Configuration struct {
	Address       string `env:"ADDRESS" envDefault:":8080"`
	AtlasUri      string `env:"ATLAS_URI,required"`
	MongoDbName   string `env:"MONGODB_NAME,required"`
	MongoCluster  string `env:"MONGO_CLUSTERS,required"`
	MongoUsername string `env:"MONGO_USERNAME,required"`
	MongoPassword string `env:"MONGO_PASSWORD,required"`
	AdminKey      string `env:"ADMIN_API_KEY,required"`
	AccessSecret  string `env:"ACCESS_SECRET,required"`
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
