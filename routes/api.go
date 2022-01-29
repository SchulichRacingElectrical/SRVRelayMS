package routes

import (
	"database-ms/app/handlers"
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
	userApi := handlers.NewUserAPI(userService)

	// Routes
	// TODO: Create middle ware for just token, just key, or both
	publicEndpoints := c.Group("/database")
	{
		// Organization
		publicEndpoints.GET("/organizations", organizationAPI.GetOrganizations)

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

	// TODO move middleware to middleware folder
	// privateEndpoints := c.Group(DbRoute, middleware.AuthorizationMiddleware())
	// {

	// TODO create authorizaion middleware
	// Organization
	// privateEndpoints.GET("/organization", controllers.GetOrganization)
	publicEndpoints.POST("/organizations", organizationAPI.Create)
	// privateEndpoints.PUT("/organization", controllers.PutOrganization)
	// privateEndpoints.DELETE("/organization", controllers.DeleteOrganization)

	// // User
	publicEndpoints.GET("/users", userApi.GetUsers)
	//publicEndpoints.GET("/users/:userId", controllers.GetUser)
	publicEndpoints.POST("/users", userApi.Create)
	// privateEndpoints.PUT("/users", controllers.PutUser)

	// }
}
