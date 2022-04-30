package routes

import (
	handlers "database-ms/app/handlers"
	middleware "database-ms/app/middleware"
	services "database-ms/app/services"
	config "database-ms/config"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

func InitializeRoutes(c *gin.Engine, mgoDbSession *mgo.Session, conf *config.Configuration) {
	// Initialize APIs
	organizationAPI := handlers.NewOrganizationAPI(services.NewOrganizationService(mgoDbSession, conf))
	userAPI := handlers.NewUserAPI(services.NewUserService(mgoDbSession, conf))
	authAPI := handlers.NewAuthAPI(services.NewUserService(mgoDbSession, conf))
	thingService := services.NewThingService(mgoDbSession, conf)
	thingAPI := handlers.NewThingAPI(thingService)
	sensorAPI := handlers.NewSensorAPI(services.NewSensorService(mgoDbSession, conf), thingService)
	operatorService := services.NewOperatorService(mgoDbSession, conf) 
	operatorAPI := handlers.NewOperatorAPI(operatorService)
	thingOperatorAPI := handlers.NewThingOperatorAPI(services.NewThingOperatorService(mgoDbSession, conf), thingService, operatorService)

	// Declare public endpoints
	publicEndpoints := c.Group("") 
	{
		organizationEndpoints := publicEndpoints.Group("/organizations") 
		{
			organizationEndpoints.GET("", organizationAPI.GetOrganizations)
			organizationEndpoints.POST("", organizationAPI.CreateOrganization)
		}
	}

	// Declare auth endpoints
	authEndpoints := c.Group("/auth")
	{
		authEndpoints.GET("/validate", authAPI.Validate)
		authEndpoints.POST("/login", authAPI.Login)
		authEndpoints.POST("/signup", authAPI.SignUp)
		authEndpoints.POST("/signout", authAPI.SignOut)
	}

	// Declare private (auth required) endpoints
	privateEndpoints := c.Group("", middleware.AuthorizationMiddleware(conf, mgoDbSession)) 
	{
		organizationEndpoints := privateEndpoints.Group("/organization")
		{
			organizationEndpoints.GET("", organizationAPI.GetOrganization)
			organizationEndpoints.PUT("", organizationAPI.UpdateOrganization)
			organizationEndpoints.PUT("/issueNewAPIKey", organizationAPI.IssueNewAPIKey)
			organizationEndpoints.DELETE("/:organizationId", organizationAPI.DeleteOrganization)
		}	

		userEndpoints := privateEndpoints.Group("/users")
		{
			userEndpoints.GET("", userAPI.GetUsers)
			userEndpoints.PUT("", userAPI.UpdateUser)
			userEndpoints.PUT("/promote", userAPI.ChangeUserRole)
			userEndpoints.DELETE("/:userId", userAPI.DeleteUser)
		}

		thingEndpoints := privateEndpoints.Group("/things")
		{
			thingEndpoints.GET("", thingAPI.GetThings)
			thingEndpoints.POST("", thingAPI.CreateThing)
			thingEndpoints.PUT("", thingAPI.UpdateThing)
			thingEndpoints.DELETE("/:thingId", thingAPI.DeleteThing)	
		}

		sensorEndpoints := privateEndpoints.Group("/sensors")
		{
			sensorEndpoints.POST("", sensorAPI.CreateSensor)
			sensorEndpoints.PUT("", sensorAPI.UpdateSensor)
			sensorEndpoints.DELETE("/:sensorId", sensorAPI.DeleteSensor)
			thingIdEndpoints := sensorEndpoints.Group("/thing/:thingId")
			{
				thingIdEndpoints.GET("", sensorAPI.FindThingSensors)
				thingIdEndpoints.GET("/lastUpdate/:lastUpdate", sensorAPI.FindUpdatedSensors)
			}	
		}

		operatorEndpoints := privateEndpoints.Group("/operators")
		{
			operatorEndpoints.POST("", operatorAPI.CreateOperator)
			operatorEndpoints.GET("", operatorAPI.GetOperators)
			operatorEndpoints.PUT("", operatorAPI.UpdateOperator)
			operatorEndpoints.DELETE("/:operatorId", operatorAPI.DeleteOperator)
		}

		thingOperatorEndpoints := privateEndpoints.Group("/thing_operator")
		{
			thingOperatorEndpoints.POST("", thingOperatorAPI.CreateThingOperatorAssociation)
			thingOperatorEndpoints.DELETE("/thing/:thingId/operator/:operatorId", thingOperatorAPI.DeleteThingOperator)
		}

		// // TODO
		// runEndpoints := privateEndpoints.Group("/runs")
		// {
		// 	runEndpoints.GET("/:thingId", )
		// 	runEndpoints.GET("/:runId/file", )
		// 	runEndpoints.GET("/:runId/comments", )
		// 	runEndpoints.POST("", )
		// 	runEndpoints.POST("/:runId/file", )
		// 	runEndpoints.POST("/comment", )
		// 	runEndpoints.PUT("", )
		// 	runEndpoints.PUT("/comment", )
		// 	runEndpoints.DELETE("/:runId", )
		// 	runEndpoints.DELETE("/comment/:commentId", )
		// }

		// // TODO
		// sessionEndpoints := privateEndpoints.Group("/sessions")
		// {
		// 	sessionEndpoints.GET("/:thingId", )
		// 	sessionEndpoints.GET("/zip", )
		// 	sessionEndpoints.POST("", )
		// 	sessionEndpoints.POST("/comment", )
		// 	sessionEndpoints.PUT("", )
		// 	sessionEndpoints.PUT("/comment", )
		// 	sessionEndpoints.DELETE("", )
		// 	sessionEndpoints.DELETE("/comment", )
		// }

		// // TODO
		// rawDataPresetEndpoints := privateEndpoints.Group("/rawdatapresets")
		// {
		// 	rawDataPresetEndpoints.GET("/:thingId", )
		// 	rawDataPresetEndpoints.POST("", )
		// 	rawDataPresetEndpoints.PUT("", )
		// 	rawDataPresetEndpoints.DELETE("/:rdpId", )
		// }

		// // TODO
		// chartPresetEndpoints := privateEndpoints.Group("/chartpresets")
		// {
		// 	chartPresetEndpoints.GET("/:thingId", )
		// 	chartPresetEndpoints.POST("", )
		// 	chartPresetEndpoints.PUT("", )
		// 	chartPresetEndpoints.DELETE("/:cpId", )
		// }
	}	
}
