package handlers

import (
	"database-ms/app/middleware"
	"database-ms/app/model"
	services "database-ms/app/services"
	utils "database-ms/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SessionHandler struct {
	session  services.SessionServiceInterface
	operator services.OperatorServiceInterface
	thing    services.ThingServiceInterface
}

func NewSessionAPI(
	sessionService services.SessionServiceInterface,
	operatorService services.OperatorServiceInterface,
	thingService services.ThingServiceInterface,
) *SessionHandler {
	return &SessionHandler{
		session:  sessionService,
		operator: operatorService,
		thing:    thingService,
	}
}

func (handler *SessionHandler) CreateSession(ctx *gin.Context) {
	var newSession model.Session
	err := ctx.BindJSON(&newSession)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Guard against non-admin requests
	if !middleware.IsAuthorizationAtLeast(ctx, "Admin") {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Guard against cross-tenant writes
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, perr := handler.thing.FindById(ctx, newSession.ThingId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingsNotFound))
		return
	}
	if organization.Id != thing.OrganizationId {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to create the session
	err = handler.session.CreateSession(ctx.Request.Context(), &newSession)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Send the response
	result := utils.SuccessPayload(newSession, "Successfully created collection.")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *SessionHandler) GetSessions(ctx *gin.Context) {
	// Attempt to read from the params
	thingId, err := uuid.Parse(ctx.Param("thingId"))
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Attempt to find the thing
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, perr := handler.thing.FindById(ctx.Request.Context(), thingId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))
		return
	}

	// Guard against cross-tenant reading
	if thing.OrganizationId != organization.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to read the sessions
	sessions, perr := handler.session.GetSessionsByThingId(ctx.Request.Context(), thingId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.SessionsNotFound))
		return
	}

	// Send the response
	result := utils.SuccessPayload(sessions, "Successfully retrieved collections")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *SessionHandler) UpdateSession(ctx *gin.Context) {
	// Attempt to extract the body
	var updatedSession model.Session
	err := ctx.BindJSON(&updatedSession)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Guard against non-admin lead or admin
	if !middleware.IsAuthorizationAtLeast(ctx, "Lead") {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to get the thing
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, perr := handler.thing.FindById(ctx, updatedSession.ThingId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))
		return
	}

	// Guard against cross-tenant updates
	if thing.OrganizationId != organization.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to update the collection
	perr = handler.session.UpdateSession(ctx.Request.Context(), &updatedSession)
	if perr != nil {
		if perr.Code == "23505" {
			utils.Response(ctx, http.StatusConflict, utils.NewHTTPError(perr.Error()))
		} else {
			utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		}
		return
	}

	// Send the response
	result := utils.SuccessPayload(nil, "Successfully updated")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *SessionHandler) DeleteSession(ctx *gin.Context) {
	// Attempt to read from the params
	sessionId, err := uuid.Parse(ctx.Param("sessionId"))
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Attempt to find the session
	organization, _ := middleware.GetOrganizationClaim(ctx)
	collection, perr := handler.session.FindById(ctx.Request.Context(), sessionId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.SensorsNotFound))
		return
	}

	// Attempt to find the thing
	thing, perr := handler.thing.FindById(ctx, collection.ThingId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))
		return
	}

	// Guard against cross-tenant deletion
	if thing.OrganizationId != organization.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to delete the session
	perr = handler.session.DeleteSession(ctx.Request.Context(), sessionId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Send the response
	result := utils.SuccessPayload(nil, "Successfully deleted")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *SessionHandler) AddComment(ctx *gin.Context) {
	var newComment model.SessionComment
	err := ctx.BindJSON(&newComment)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Guard against cross-tenant write
	organization, _ := middleware.GetOrganizationClaim(ctx)
	collection, perr := handler.session.FindById(ctx, newComment.SessionId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.CollectionNotFound))
		return
	}
	thing, perr := handler.thing.FindById(ctx, collection.ThingId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))
		return
	}
	if thing.OrganizationId != organization.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	newComment.LastUpdate = utils.CurrentTimeInMilli()

	// Attempt to create the collection
	err = handler.session.AddComment(ctx.Request.Context(), &newComment)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Send the response
	result := utils.SuccessPayload(newComment, "Successfully added comment.")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *SessionHandler) GetComments(ctx *gin.Context) {
	// Attempt to read from the params
	sessionId, err := uuid.Parse(ctx.Param("sessionId"))
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Guard against cross-tenant read
	organization, _ := middleware.GetOrganizationClaim(ctx)
	session, perr := handler.session.FindById(ctx, sessionId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.CollectionNotFound))
		return
	}
	thing, perr := handler.thing.FindById(ctx, session.ThingId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))
		return
	}
	if thing.OrganizationId != organization.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to read the comments
	comments, err := handler.session.GetComments(ctx.Request.Context(), sessionId)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.CommentsNotFound))
		return
	}

	// Send the response
	result := utils.SuccessPayload(comments, "Successfully retrieved comments")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *SessionHandler) UpdateCommentContent(ctx *gin.Context) {
	var updatedComment model.SessionComment
	err := ctx.BindJSON(&updatedComment)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Guard against cross-user update
	user, err := middleware.GetUserClaim(ctx)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}
	updatedComment.UserId = user.Id
	comment, err := handler.session.GetComment(ctx, updatedComment.Id)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.CommentNotFound))
		return
	}
	if comment.UserId != user.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	updatedComment.LastUpdate = utils.CurrentTimeInMilli()

	// Attempt to create the collection
	err = handler.session.UpdateCommentContent(ctx.Request.Context(), &updatedComment)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Send the response
	result := utils.SuccessPayload(nil, "Successfully updated")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *SessionHandler) DeleteComment(ctx *gin.Context) {
	// Attempt to read from the params
	commentId, err := uuid.Parse(ctx.Param("commentId"))
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Guard against cross-user update
	user, err := middleware.GetUserClaim(ctx)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}
	comment, err := handler.session.GetComment(ctx, commentId)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.CommentNotFound))
		return
	}
	if comment.UserId != user.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to delete the comment
	err = handler.session.DeleteComment(ctx.Request.Context(), commentId)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Send the response
	result := utils.SuccessPayload(nil, "Successfully deleted")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *SessionHandler) UploadFile(c *gin.Context) {
	// // Check if run exist
	// run, err := handler.run.FindById(c.Request.Context(), c.PostForm("runId"))
	// if err != nil {
	// 	utils.Response(c, http.StatusNotFound, utils.NewHTTPError(utils.RunNotFound))
	// 	return
	// }

	// // Check if operator exist
	// operator, err := handler.operator.FindById(c.Request.Context(), c.PostForm("operatorId"))
	// if err != nil {
	// 	utils.Response(c, http.StatusNotFound, utils.NewHTTPError(utils.OperatorNotFound))
	// 	return
	// }

	// // Check if thing exist
	// thing, err := handler.thing.FindById(c.Request.Context(), c.PostForm("thingId"))
	// if err != nil {
	// 	utils.Response(c, http.StatusNotFound, utils.NewHTTPError(utils.ThingNotFound))
	// 	return
	// }

	// // Check if run alread has a file
	// if runMetadata, _ := handler.run.GetRunFileMetaData(c.Request.Context(), c.PostForm("runId")); runMetadata != nil {
	// 	utils.Response(c, http.StatusNotFound, utils.NewHTTPError(utils.RunHasAssociatedFile))
	// 	return
	// }

	// runFileMetadata := models.RunFileUpload{
	// 	OperatorId:      operator.ID,
	// 	RunId:           run.ID,
	// 	ThingID:         thing.ID,
	// 	UploadDateEpoch: utils.CurrentTimeInMilli(),
	// }

	// file, err := c.FormFile("file")
	// if err != nil {
	// 	result := utils.NewHTTPError(utils.NoFileReceived)
	// 	utils.Response(c, http.StatusBadRequest, result)
	// 	return
	// }

	// // Verify file extension
	// if extension := filepath.Ext(file.Filename); extension != ".csv" {
	// 	fmt.Println(extension)
	// 	result := utils.NewHTTPError(utils.NotCsv)
	// 	utils.Response(c, http.StatusBadRequest, result)
	// 	return
	// }

	// // Save file
	// err = handler.run.UploadFile(c.Request.Context(), &runFileMetadata, file)
	// if err != nil {
	// 	utils.Response(c, http.StatusInternalServerError, utils.NewHTTPError(utils.FileNotUploaded))
	// }

	// result := utils.SuccessPayload(nil, "Successfully uploaded file")
	// utils.Response(c, http.StatusOK, result)
}

func (handler *SessionHandler) DownloadFile(c *gin.Context) {
	// Check if run alread has a file
	// runMetadata, err := handler.run.GetRunFileMetaData(c.Request.Context(), c.PostForm("runId"))
	// if err != nil {
	// 	utils.Response(c, http.StatusNotFound, utils.NewHTTPError(utils.RunHasNoAssociatedFile))
	// 	return
	// }

	// byteFile, err := handler.run.DownloadFile(c.Request.Context(), c.PostForm("runId"))
	// if err != nil {
	// 	utils.Response(c, http.StatusInternalServerError, utils.NewHTTPError(utils.CannotRetrieveFile))
	// 	return
	// }
	// c.Header("Content-Disposition", "attachment; filename="+runMetadata.FileName)
	// c.Data(http.StatusOK, "application/octet-stream", byteFile)
}

func (handler *SessionHandler) GetDatumBySessionIdAndSensorId(c *gin.Context) {
	// datumArray, err := handler.datum.FindBySessionIdAndSensorId(c.Request.Context(), c.Param("sessionId"), c.Param("sensorId"))
	// if err != nil {
	// 	utils.Response(c, http.StatusBadRequest, utils.NewHTTPError(utils.DatumNotFound))
	// } else {
	// 	result := utils.SuccessPayload(datumArray, "Successfully retrieved datum")
	// 	utils.Response(c, http.StatusOK, result)
	// }
}
