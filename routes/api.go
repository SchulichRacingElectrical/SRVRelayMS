package routes

import (
	"database-ms/app/handlers"
	sensorRepo "database-ms/app/repositories/sensor"
	sensorSrv "database-ms/app/services/sensor"
	"database-ms/config"
	"database-ms/controllers"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

var (
	DbRoute = "/database"
)

func InitializeRoutes(c *gin.Engine, dbSession *mgo.Session, conf *config.Configuration) {
	sensorRepository := sensorRepo.New(dbSession, conf)
	sensorService := sensorSrv.New(sensorRepository)
	sensorAPI := handlers.NewSensorAPI(sensorService)

	// Routes
	// TODO: Create middle ware for just token, just key, or both
	publicEndpoints := c.Group("/database")
	{
		// Organization
		publicEndpoints.GET("/organizations", controllers.GetOrganizations)

		// Sensor
		// TODO move later as private endpoint
		publicEndpoints.POST("/sensors", sensorAPI.Create)
	}

	privateEndpoints := c.Group(DbRoute, controllers.AuthorizationMiddleware())
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

	}
}
