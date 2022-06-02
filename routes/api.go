package routes

import (
	handlers "database-ms/app/handlers"
	middleware "database-ms/app/middleware"
	services "database-ms/app/services"
	config "database-ms/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitializeRoutes(c *gin.Engine, db *gorm.DB, conf *config.Configuration) {
	// Initialize APIs
	organizationService := services.NewOrganizationService(db, conf)
	organizationAPI := handlers.NewOrganizationAPI(organizationService)
	userAPI := handlers.NewUserAPI(services.NewUserService(db, conf))
	authAPI := handlers.NewAuthAPI(services.NewUserService(db, conf), organizationService)
	thingService := services.NewThingService(db, conf)
	thingAPI := handlers.NewThingAPI(thingService)
	sensorAPI := handlers.NewSensorAPI(services.NewSensorService(db, conf), thingService)
	operatorService := services.NewOperatorService(db, conf)
	operatorAPI := handlers.NewOperatorAPI(operatorService)
	sessionService := services.NewSessionService(db, conf)
	sessionAPI := handlers.NewSessionAPI(sessionService, thingService, conf.FilePath)
	collectionService := services.NewCollectionService(db, conf)
	collectionAPI := handlers.NewCollectionAPI(collectionService, thingService)
	commentAPI := handlers.NewCommentAPI(services.NewCommentService(db, conf), thingService, sessionService, collectionService)
	rawDataPresetAPI := handlers.NewRawDataPresetAPI(services.NewRawDataPresetService(db, conf), thingService)
	chartPresetAPI := handlers.NewChartPresetAPI(services.NewChartPresetService(db, conf), thingService)

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
	privateEndpoints := c.Group("", middleware.AuthorizationMiddleware(conf, db))
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

		sessionEndpoints := privateEndpoints.Group("/sessions")
		{
			sessionEndpoints.POST("", sessionAPI.CreateSession)
			sessionEndpoints.GET("/thing/:thingId", sessionAPI.GetSessions)
			sessionEndpoints.PUT("", sessionAPI.UpdateSession)
			sessionEndpoints.DELETE("/:sessionId", sessionAPI.DeleteSession)
			sessionEndpoints.POST("/:sessionId/file", sessionAPI.UploadFile)
			sessionEndpoints.GET("/:sessionId/file", sessionAPI.DownloadFile)
		}

		collectionEndpoints := privateEndpoints.Group("/collections")
		{
			collectionEndpoints.POST("", collectionAPI.CreateCollection)
			collectionEndpoints.GET("/thing/:thingId", collectionAPI.GetCollections)
			collectionEndpoints.PUT("", collectionAPI.UpdateCollections)
			collectionEndpoints.DELETE("/:collectionId", collectionAPI.DeleteCollection)
		}

		commentEndpoints := privateEndpoints.Group("/comments")
		{
			commentEndpoints.POST("", commentAPI.CreateComment)
			commentEndpoints.GET("/:contextId", commentAPI.GetComments)
			commentEndpoints.PUT("", commentAPI.UpdateComment)
			commentEndpoints.DELETE("/:commentId", commentAPI.DeleteComment)
		}

		rawDataPresetEndpoints := privateEndpoints.Group("/rawDataPreset")
		{
			rawDataPresetEndpoints.GET("/thing/:thingId", rawDataPresetAPI.GetRawDataPresets)
			rawDataPresetEndpoints.POST("", rawDataPresetAPI.CreateRawDataPreset)
			rawDataPresetEndpoints.PUT("", rawDataPresetAPI.UpdateRawDataPreset)
			rawDataPresetEndpoints.DELETE("/:rawDataPresetId", rawDataPresetAPI.DeleteRawDataPreset)
		}

		chartPresetEndpoints := privateEndpoints.Group("/chartPreset")
		{
			chartPresetEndpoints.GET("/thing/:thingId", chartPresetAPI.GetChartPresets)
			chartPresetEndpoints.POST("", chartPresetAPI.CreateChartPreset)
			chartPresetEndpoints.PUT("", chartPresetAPI.UpdateChartPreset)
			chartPresetEndpoints.DELETE("/:chartPresetId", chartPresetAPI.DeleteChartPreset)
		}

		dataEndpoints := privateEndpoints.Group("/data")
		{
			dataEndpoints.GET("/:sessionId/:sensorId", sessionAPI.GetDatumBySessionIdAndSensorId)
		}
	}
}
