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
	//TODO: set config file for server info

	//Initialiaze firebase database
	err := databases.Database.Init()
	if err != nil {
		return err
	}

	//TODO: Set Gin logger

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

	databaseHandlers := m.router.Group("/database")
	{
		// Organization
		databaseHandlers.GET("/organizations", controllers.GetOrganizations)
		databaseHandlers.GET("/organizations/:organizationId", controllers.GetOrganization)
		databaseHandlers.POST("/organizations", controllers.PostOrganization)

		// User
		databaseHandlers.GET("/users/:organizationId", controllers.GetUsers)
		databaseHandlers.GET("/users/:organizationId/:userId", controllers.GetUser)
		databaseHandlers.POST("/users", controllers.PostUser)
		databaseHandlers.PUT("/users", controllers.PutUser)

		// Sensor
		databaseHandlers.GET("/sensors", controllers.GetSensors)
		databaseHandlers.GET("/sensors/:sid", controllers.GetSensor)
		databaseHandlers.POST("/sensors", controllers.PostSensor)
		databaseHandlers.DELETE("/sensors/:sid", controllers.DeleteSensor)
	}

	m.router.Run(":8080")
}
