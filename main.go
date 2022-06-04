package main

import (
	"database-ms/app/databases"
	"database-ms/app/subscriber"
	"database-ms/config"
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
	InitializeRoutes(router, db, conf)
	router.Use(cors.Default())
	gin.SetMode(gin.ReleaseMode)

	// Redis IoT Sub
	subscriber.Initialize(conf, db)

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
