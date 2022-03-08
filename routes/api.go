package routes

import (
	"database-ms/app/handlers"
	"database-ms/app/middleware"
	organizationSrv "database-ms/app/services/organization"
	sensorSrv "database-ms/app/services/sensor"
	thingSrv "database-ms/app/services/thing"
	userSrv "database-ms/app/services/user"
	"database-ms/config"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

var (
	DbRoute = "/database"
)

func InitializeRoutes(c *gin.Engine, mgoDbSession *mgo.Session, conf *config.Configuration) {
	sensorService := sensorSrv.NewSensorService(mgoDbSession, conf)
	sensorAPI := handlers.NewSensorAPI(sensorService)

	organizationService := organizationSrv.NewOrganizationService(mgoDbSession, conf)
	organizationAPI := handlers.NewOrganizationAPI(organizationService)

	userService := userSrv.NewUserService(mgoDbSession, conf)
	userAPI := handlers.NewUserAPI(userService)

	thingService := thingSrv.NewThingService(mgoDbSession, conf)
	thingAPI := handlers.NewThingAPI(thingService)

	// Routes
	// TODO move later as private endpoint
	publicEndpoints := c.Group("/database")
	{
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

	authEndpoints := c.Group("/auth")
	{
		// Authentication
		authEndpoints.POST("/login", userAPI.Login)
		authEndpoints.POST("/signup", userAPI.Create)
	}

	// Temp endpoints
	privateEndpoints := c.Group(DbRoute, middleware.AuthorizationMiddleware(conf, mgoDbSession))
	{
		// Organizations
		privateEndpoints.POST("/organizations", organizationAPI.Create)

		// Users
		privateEndpoints.GET("/users/userId/:userId", userAPI.GetUser)
	}
}
