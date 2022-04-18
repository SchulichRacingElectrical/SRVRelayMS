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
	DbRoute = ""
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
	publicEndpoints := c.Group("")
	{
		// Most of these need to be private endpoints
		// Sensor
		sensorsEndpoints := publicEndpoints.Group("/sensors")
		{
			sensorsEndpoints.POST("", sensorAPI.Create)
			// No need for /sensorId, just use /:sensorId, do we even need this endpoint?
			sensorsEndpoints.GET("/:sensorId", sensorAPI.FindBySensorId)
			// Don't need any route at all, the sensor Id will be in the body
			sensorsEndpoints.PUT("/:sensorId", sensorAPI.Update)
			// Don't need /sensorId/:sensorId at all, the sensor id will be in the body
			sensorsEndpoints.DELETE("/:sensorId", sensorAPI.Delete)

			thingIdEndpoints := sensorsEndpoints.Group("/thing/sensors")
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
				thingIdEndpoints.GET("", thingAPI.GetThings)
				thingIdEndpoints.PUT("", thingAPI.UpdateThing)
				thingIdEndpoints.DELETE("", thingAPI.Delete)
			}
		}

		// Organizations
		organizationEndpoints := publicEndpoints.Group("/organizations")
		{
			organizationEndpoints.GET("", organizationAPI.FindAllOrganizations)
		}
	}

	authEndpoints := c.Group("/auth")
	{
		authEndpoints.POST("/login", userAPI.Login)
		authEndpoints.POST("/signup", userAPI.Create)
	}

	privateEndpoints := c.Group(DbRoute, middleware.AuthorizationMiddleware(conf, mgoDbSession))
	{
		// Organizations
		organizationEndpoints := privateEndpoints.Group("/organizations")
		{
			organizationEndpoints.POST("", organizationAPI.Create)
			organizationEndpoints.GET("organizationId/:organizationId", organizationAPI.FindByOrganizationId)
		}

		// Users
		userEndpoints := privateEndpoints.Group("/users")
		{
			// Don't use /userId, just /:userId
			userEndpoints.GET("/userId/:userId", userAPI.GetUser)
		}
	}
}
