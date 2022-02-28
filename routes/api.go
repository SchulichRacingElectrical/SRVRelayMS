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

	authEndpoints := c.Group("/auth")
	{
		// Authentication
		authEndpoints.POST("/login", userAPI.Login)
		authEndpoints.POST("/signup", userAPI.Create)
	}

	// Temp endpoints
	privateEndpoints := c.Group(DbRoute, middleware.AuthorizationMiddleware(conf, dbSession))
	{
		// Organizations
		privateEndpoints.POST("/organizations", organizationAPI.Create)

		// Users
		privateEndpoints.GET("/users/userId/:userId", userAPI.GetUser)
	}
}
