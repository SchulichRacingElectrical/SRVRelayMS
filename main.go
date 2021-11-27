package main

import (
	"database-ms/config"
	"database-ms/controllers"
	"database-ms/databases"

	"github.com/gin-gonic/gin"
)

type Main struct {
	router *gin.Engine
}

func (m *Main) initServer(conf *config.Configuration) {

	// Initialiaze firebase database
	// err := databases.Firebase.Init()
	// if err != nil {
	// 	return err
	// }

	// Initialize MongoDB
	databases.Mongo.Init(conf.AtlasUri, conf.MongoDbName)

	// TODO set Gin logger

	m.router = gin.Default()
}

func main() {
	m := Main{}

	conf := config.NewConfig("./env")

	// Initialize server
	m.initServer(conf)

	// defer databases.Firebase.Close()
	defer databases.Mongo.Close()

	// TODO: Create middle ware for just token, just key, or both
	publicEndpoints := m.router.Group("/database")
	{
		// Organization
		publicEndpoints.GET("/organizations", controllers.GetOrganizations)
	}

	privateEndpoints := m.router.Group("/database", controllers.AuthorizationMiddleware())
	{
		// Organization
		privateEndpoints.GET("/organization", controllers.GetOrganization)
		privateEndpoints.POST("/organization", controllers.PostOrganization)
		privateEndpoints.PUT("/organization", controllers.PutOrganization)
		privateEndpoints.DELETE("/organization", controllers.DeleteOrganization)

		// User
		privateEndpoints.GET("/users", controllers.GetUsers)
		privateEndpoints.GET("/users/:userId", controllers.GetUser)
		privateEndpoints.POST("/users", controllers.PostUser)
		privateEndpoints.PUT("/users", controllers.PutUser)

		// Sensor
		privateEndpoints.GET("/sensors", controllers.GetSensors)
		privateEndpoints.GET("/sensors/:sid", controllers.GetSensor)
		privateEndpoints.POST("/sensors", controllers.CreateSensor)
		privateEndpoints.DELETE("/sensors/:sid", controllers.DeleteSensor)
	}

	m.router.Run(":8080")
}
