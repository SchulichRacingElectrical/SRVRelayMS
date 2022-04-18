package routes

import (
	"database-ms/app/handlers"
	"database-ms/app/middleware"
	services "database-ms/app/services"
	"database-ms/config"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

func InitializeRoutes(c *gin.Engine, mgoDbSession *mgo.Session, conf *config.Configuration) {
	// Initialize APIs
	organizationAPI := handlers.NewOrganizationAPI(services.NewOrganizationService(mgoDbSession, conf))
	userAPI := handlers.NewUserAPI(services.NewUserService(mgoDbSession, conf))
	thingAPI := handlers.NewThingAPI(services.NewThingService(mgoDbSession, conf))
	sensorAPI := handlers.NewSensorAPI(services.NewSensorService(mgoDbSession, conf))

	// Declare public endpoints
	publicEndpoints := c.Group("") 
	{
		organizationEndpoints := publicEndpoints.Group("/organizations") 
		{
			organizationEndpoints.GET("", organizationAPI.GetOrganizations)
			organizationEndpoints.POST("", organizationAPI.Create)
		}
	}

	// Declare auth endpoints
	authEndpoints := c.Group("/auth")
	{
		authEndpoints.POST("/login", userAPI.Login)
		authEndpoints.POST("/signup", userAPI.Create)
	}

	// Declare private (auth required) endpoints
	privateEndpoints := c.Group("", middleware.AuthorizationMiddleware(conf, mgoDbSession)) 
	{
		organizationEndpoints := privateEndpoints.Group("/organizations")
		{
			organizationEndpoints.GET("organizationId/:organizationId", organizationAPI.GetOrganization)
		}	

		userEndpoints := privateEndpoints.Group("/users")
		{
			// Don't use /userId, just /:userId
			userEndpoints.GET("/userId/:userId", userAPI.GetUser)
		}

		thingEndpoints := privateEndpoints.Group("/things")
		{
			thingEndpoints.POST("", thingAPI.Create)	
			thingIdEndpoints := thingEndpoints.Group("/:thingId")
			{
				thingIdEndpoints.GET("", thingAPI.GetThings)
				thingIdEndpoints.PUT("", thingAPI.UpdateThing)
				thingIdEndpoints.DELETE("", thingAPI.Delete)	
			}
		}

		sensorEndpoints := privateEndpoints.Group("/sensors")
		{
			sensorEndpoints.POST("", sensorAPI.Create)
			sensorEndpoints.GET("/:sensorId", sensorAPI.FindBySensorId) // Do we need this?
			sensorEndpoints.PUT("", sensorAPI.Update)
			sensorEndpoints.DELETE("/:sensorId", sensorAPI.Delete)

			thingIdEndpoints := sensorEndpoints.Group("/thing/sensors")
			{
				thingIdEndpoints.GET("/:thingId", sensorAPI.FindThingSensors)
				thingIdEndpoints.GET("/:thingId/lastUpdate/:lastUpdate", sensorAPI.FindUpdatedSensor)
			}	
		}
	}	
}
