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

	// Attempt to delete the collection
	perr = handler.session.DeleteSession(ctx.Request.Context(), sessionId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Send the response
	result := utils.SuccessPayload(nil, "Successfully deleted")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *SessionHandler) AddComment(c *gin.Context) {
	// var comment models.Comment
	// c.BindJSON(&comment)

	// _, err := handler.run.FindById(c.Request.Context(), c.Param("runId"))
	// if err == nil {
	// 	err := handler.comment.AddComment(c.Request.Context(), utils.Run, c.Param("runId"), &comment)
	// 	if err == nil {
	// 		result := utils.SuccessPayload(nil, "Successfully added comment.")
	// 		utils.Response(c, http.StatusOK, result)
	// 	} else {
	// 		utils.Response(c, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
	// 	}
	// } else {
	// 	utils.Response(c, http.StatusNotFound, utils.NewHTTPError(utils.RunNotFound))
	// }
}

func (handler *SessionHandler) GetComments(c *gin.Context) {
	// comments, err := handler.comment.GetComments(c.Request.Context(), utils.Run, c.Param("runId"))
	// if err == nil {
	// 	result := utils.SuccessPayload(comments, "Successfully retrieved comments")
	// 	utils.Response(c, http.StatusOK, result)
	// } else {
	// 	utils.Response(c, http.StatusBadRequest, utils.NewHTTPError(utils.CommentsNotFound))
	// }
}

func (handler *SessionHandler) UpdateCommentContent(c *gin.Context) {
	// var updatedComment models.Comment
	// c.BindJSON(&updatedComment)

	// err := handler.comment.UpdateCommentContent(c.Request.Context(), c.Param("commentId"), &updatedComment)
	// if err == nil {
	// 	result := utils.SuccessPayload(nil, "Successfully updated comment")
	// 	utils.Response(c, http.StatusOK, result)
	// } else {
	// 	var errMsg string
	// 	switch err.Error() {
	// 	case utils.CommentDoesNotExist, utils.CommentCannotUpdateOtherUserComment:
	// 		errMsg = err.Error()
	// 	default:
	// 		errMsg = utils.BadRequest
	// 	}
	// 	result := utils.NewHTTPError(errMsg)
	// 	utils.Response(c, http.StatusBadRequest, result)
	// }
}

func (handler *SessionHandler) DeleteComment(c *gin.Context) {
	// var requestBody models.Comment
	// c.BindJSON(&requestBody)

	// if !requestBody.UserID.IsZero() {
	// 	err := handler.comment.DeleteComment(c.Request.Context(), utils.Run, c.Param("commentId"), requestBody.UserID.Hex())
	// 	if err == nil {
	// 		result := utils.SuccessPayload(nil, "Successfully deleted comment")
	// 		utils.Response(c, http.StatusOK, result)
	// 	} else {
	// 		var errMsg string
	// 		switch err.Error() {
	// 		case utils.CommentDoesNotExist, utils.CommentCannotUpdateOtherUserComment:
	// 			errMsg = err.Error()
	// 		default:
	// 			errMsg = utils.BadRequest
	// 		}
	// 		result := utils.NewHTTPError(errMsg)
	// 		utils.Response(c, http.StatusBadRequest, result)
	// 	}
	// } else {
	// 	result := utils.NewHTTPError(utils.UserIdMissing)
	// 	utils.Response(c, http.StatusBadRequest, result)
	// }
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
