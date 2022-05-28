package main

import (
	"database-ms/config"
	"database-ms/databases"
	"database-ms/redisHandler"
	"database-ms/routes"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

func main() {
	// Initialize config
	conf := config.NewConfig("./env")

	// Connect to the Postgres DB
	db := databases.InitPostgres(conf)
	// defer db.Close() -> Need to close somehow?

	// Router
	router := gin.Default()
	routes.InitializeRoutes(router, db, conf)
	router.Use(cors.Default())

	// Redis IoT Sub
	redisHandler.Initialize(conf, db)

	// Server config
	srv := &http.Server{
		Handler:      router,
		Addr:         conf.Address,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	// Serving microservice at specified port
	log.Println("SRV-DB-MS running at ", conf.Address)
	log.Fatal(srv.ListenAndServe())
}
