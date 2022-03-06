package routes

import (
	"database-ms/app/handlers"
	"database-ms/app/middleware"
	sensorSrv "database-ms/app/services/sensor"
	thingSrv "database-ms/app/services/thing"
	"database-ms/config"
	"database-ms/controllers"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

var (
	DbRoute = "/database"
)

func InitializeRoutes(c *gin.Engine, mgoDbSession *mgo.Session, conf *config.Configuration) {
	sensorService := sensorSrv.NewSensorService(mgoDbSession, conf)
	sensorAPI := handlers.NewSensorAPI(sensorService)

	thingService := thingSrv.NewThingService(mgoDbSession, conf)
	thingAPI := handlers.NewThingAPI(thingService)

	// Routes
	// TODO move later as private endpoint
	publicEndpoints := c.Group("/database")
	{
		// Organization
		publicEndpoints.GET("/organizations", controllers.GetOrganizations)

		// Sensor
		sensorsEndpoints := publicEndpoints.Group("/sensors")
		{
			sensorsEndpoints.POST("", sensorAPI.Create)
			sensorsEndpoints.GET("/sensorId/:sensorId", sensorAPI.FindBySensorId)
			sensorsEndpoints.PUT("/sensorId/:sensorId", sensorAPI.Update)
			sensorsEndpoints.DELETE("/sensorId/:sensorId", sensorAPI.Delete)

			thingIdEndpoints := sensorsEndpoints.Group("/thingId")
			{
				thingIdEndpoints.GET("/:thingId", sensorAPI.FindThingSensors)
				thingIdEndpoints.GET("/:thingId/lastUpdate/:lastUpdate", sensorAPI.FindUpdatedSensor)
			}
		}

		// Thing
		thingEndpoints := publicEndpoints.Group("/thing")
		{
			thingEndpoints.POST("", thingAPI.Create)
			thingIdEndpoints := thingEndpoints.Group("/:thingId")
			{
				thingIdEndpoints.GET("", thingAPI.GetThing)
				thingIdEndpoints.PUT("", thingAPI.UpdateThing)
				thingIdEndpoints.DELETE("", thingAPI.Delete)
			}
		}
	}

	// TODO move middleware to middleware folder
	privateEndpoints := c.Group(DbRoute, middleware.AuthorizationMiddleware())
	{

		// TODO refactor these endpoints to use multitier pattern

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
