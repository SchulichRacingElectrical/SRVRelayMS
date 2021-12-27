package main

import (
	"database-ms/config"
	"database-ms/databases"
	"database-ms/routes"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"gopkg.in/mgo.v2"
)

func main() {

	// Initialize config
	conf := config.NewConfig("./env")

	// Connect to DB
	mongoConn := databases.GetInstance(conf)

	dbSession := mongoConn.Copy()
	dbSession.SetSafe(&mgo.Safe{})
	defer dbSession.Close()

	// Router
	router := gin.Default()
	routes.InitializeRoutes(router, dbSession, conf)
	router.Use(cors.Default())

	// TODO setup Swagger
	// TODO setup logging

	// Server config
	srv := &http.Server{
		Handler:      router,
		Addr:         conf.Address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("SRV-DB-MS running at ", conf.Address)

	// Serving microservice at specified port
	log.Fatal(srv.ListenAndServe())
}
