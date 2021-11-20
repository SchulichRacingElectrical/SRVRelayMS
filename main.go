package main

import (
	"database-ms/controllers"
	"database-ms/databases"

	"github.com/gin-gonic/gin"
)

type Main struct {
	router *gin.Engine
}

func (m *Main) initServer() error {
	// TODO: set config file for server info

	// Initialiaze firebase database
	err := databases.Database.Init()
	if err != nil {
		return err
	}

	// TODO: Set Gin logger

	m.router = gin.Default()

	return nil
}

func main() {
	m := Main{}

	// Initialize server
	if m.initServer() != nil {
		return
	}

	defer databases.Database.Close()

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
		privateEndpoints.POST("/sensors", controllers.PostSensor)
		privateEndpoints.DELETE("/sensors/:sid", controllers.DeleteSensor)
	}

	m.router.Run(":8080")
}
