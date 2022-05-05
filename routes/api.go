package routes

import (
	"database-ms/app/handlers"
	"database-ms/app/middleware"
	organizationSrv "database-ms/app/services/organization"
	runSrv "database-ms/app/services/run"
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

	runService := runSrv.NewRunService(conf)
	runAPI := handlers.NewRunAPI(runService)

	// Routes
	publicEndpoints := c.Group("")
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

		// Run
		runEndpoints := publicEndpoints.Group("/run")
		{
			runEndpoints.POST("", runAPI.Create)
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
			userEndpoints.GET("/userId/:userId", userAPI.GetUser)
		}
	}
}
