package routes

import (
	"database-ms/app/handlers"
	"database-ms/app/middleware"
	organizationSrv "database-ms/app/services/organization"
	sensorSrv "database-ms/app/services/sensor"
	userSrv "database-ms/app/services/user"
	"database-ms/config"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

var (
	DbRoute = "/database"
)

func InitializeRoutes(c *gin.Engine, dbSession *mgo.Session, conf *config.Configuration) {
	sensorService := sensorSrv.NewSensorService(dbSession, conf)
	sensorAPI := handlers.NewSensorAPI(sensorService)

	organizationService := organizationSrv.NewOrganizationService(dbSession, conf)
	organizationAPI := handlers.NewOrganizationAPI(organizationService)

	userService := userSrv.NewUserService(dbSession, conf)
	userAPI := handlers.NewUserAPI(userService)

	// Routes
	// TODO: Create middle ware for just token, just key, or both
	publicEndpoints := c.Group("/database")
	{
		// Sensor
		// TODO move later as private endpoint
		publicEndpoints.POST("/sensors", sensorAPI.Create)
		publicEndpoints.GET("/sensors/sensorId/:sensorId", sensorAPI.FindBySensorId)
		thingIdEndpoints := publicEndpoints.Group("/sensors/thingId")
		{
			thingIdEndpoints.GET("/:thingId", sensorAPI.FindThingSensors)
			thingIdEndpoints.GET("/:thingId/lastUpdate/:lastUpdate", sensorAPI.FindUpdatedSensor)
		}
		publicEndpoints.PUT("/sensors/sensorId/:sensorId", sensorAPI.Update)
		publicEndpoints.DELETE("/sensors/sensorId/:sensorId", sensorAPI.Delete)
	}

	// make an api key first you dingus!!!

	// Temp endpoints
	privateEndpoints := c.Group(DbRoute, middleware.AuthorizationMiddleware(conf))
	{
		// Organizations
		privateEndpoints.POST("/organizations", organizationAPI.Create)

		// Users
		privateEndpoints.POST("/users", userAPI.Create)
	}

	// TODO move middleware to middleware folder
	// privateEndpoints := c.Group(DbRoute, middleware.AuthorizationMiddleware())
	// {
	// 	// TODO refactor these endpoints to use multitier pattern

	// 	// Organization
	// 	privateEndpoints.GET("/organization", controllers.GetOrganization)
	// 	privateEndpoints.POST("/organization", controllers.PostOrganization)
	// 	privateEndpoints.PUT("/organization", controllers.PutOrganization)
	// 	privateEndpoints.DELETE("/organization", controllers.DeleteOrganization)

	// 	// User
	// 	privateEndpoints.GET("/users", controllers.GetUsers)
	// 	privateEndpoints.GET("/users/:userId", controllers.GetUser)
	// 	privateEndpoints.POST("/users", controllers.PostUser)
	// 	privateEndpoints.PUT("/users", controllers.PutUser)

	// }
}
