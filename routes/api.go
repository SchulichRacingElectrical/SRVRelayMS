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
	organizationService := services.NewOrganizationService(mgoDbSession, conf)
	organizationAPI := handlers.NewOrganizationAPI(organizationService)
	userAPI := handlers.NewUserAPI(services.NewUserService(mgoDbSession, conf))
	authAPI := handlers.NewAuthAPI(services.NewUserService(mgoDbSession, conf), organizationService)
	thingService := services.NewThingService(mgoDbSession, conf)
	thingAPI := handlers.NewThingAPI(thingService)
	sensorAPI := handlers.NewSensorAPI(services.NewSensorService(mgoDbSession, conf), thingService)
	operatorService := services.NewOperatorService(mgoDbSession, conf)
	operatorAPI := handlers.NewOperatorAPI(operatorService)
	commentService := services.NewCommentService(conf)
	runAPI := handlers.NewRunAPI(services.NewRunService(conf), commentService, operatorService, thingService)
	sessionAPI := handlers.NewSessionAPI(services.NewSessionService(conf), commentService)

	chartPresetAPI := handlers.NewChartPresetAPI(services.NewChartPresetService(mgoDbSession, conf), thingService)
	rawDataPresetAPI := handlers.NewRawDataPresetAPI(services.NewRawDataPresetService(mgoDbSession, conf), thingService)

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

		runEndpoints := privateEndpoints.Group("/runs")
		{
			runEndpoints.POST("", runAPI.CreateRun)
			runEndpoints.GET("/thing/:thingId", runAPI.GetRuns)
			runEndpoints.PUT("", runAPI.UpdateRun)
			runEndpoints.DELETE("/:runId", runAPI.DeleteRun)

			runEndpoints.POST("/:runId/file", runAPI.UploadFile)
			runEndpoints.GET("/:runId/file", runAPI.DownloadFile)

			runEndpoints.POST("/:runId/comment", runAPI.AddComment)
			runEndpoints.GET("/:runId/comments", runAPI.GetComments)
			runEndpoints.PUT("/comment/:commentId", runAPI.UpdateCommentContent)
			runEndpoints.DELETE("/comment/:commentId", runAPI.DeleteComment)
		}

		sessionEndpoints := privateEndpoints.Group("/sessions")
		{
			sessionEndpoints.POST("", sessionAPI.CreateSession)
			sessionEndpoints.GET("/thing/:thingId", sessionAPI.GetSessions)
			sessionEndpoints.PUT("", sessionAPI.UpdateSession)
			sessionEndpoints.DELETE("/:sessionId", sessionAPI.DeleteSession)

			// TODO
			// sessionEndpoints.GET("/zip", )

			sessionEndpoints.POST("/:sessionId/comment", sessionAPI.AddComment)
			sessionEndpoints.GET("/:sessionId/comments", sessionAPI.GetComments)
			sessionEndpoints.PUT("/comment/:commentId", sessionAPI.UpdateCommentContent)
			sessionEndpoints.DELETE("/comment/:commentId", sessionAPI.DeleteComment)
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
	}
}
