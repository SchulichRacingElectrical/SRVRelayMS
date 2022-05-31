package handlers

import (
	"bytes"
	"database-ms/app/middleware"
	"database-ms/app/model"
	services "database-ms/app/services"
	utils "database-ms/utils"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SessionHandler struct {
	session  services.SessionServiceInterface
	thing    services.ThingServiceInterface
	filepath string
}

func NewSessionAPI(
	sessionService services.SessionServiceInterface,
	thingService services.ThingServiceInterface,
	filepath string,
) *SessionHandler {
	return &SessionHandler{
		session:  sessionService,
		thing:    thingService,
		filepath: filepath,
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
	session, perr := handler.session.FindById(ctx, newComment.SessionId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.SessionNotFound))
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

	newComment.LastUpdate = utils.CurrentTimeInMilli()

	// Attempt to create the session
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
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.SessionNotFound))
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

	// Attempt to create the session
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

func (handler *SessionHandler) UploadFile(ctx *gin.Context) {
	// Guard against non-admin lead or admin
	if !middleware.IsAuthorizationAtLeast(ctx, "Lead") {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to read from the params
	sessionId, err := uuid.Parse(ctx.Param("sessionId"))
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Check if session exists
	session, perr := handler.session.FindById(ctx, sessionId)
	if perr != nil {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.SessionNotFound))
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		result := utils.NewHTTPError(utils.NoFileRcvd)
		utils.Response(ctx, http.StatusBadRequest, result)
		return
	}

	// // Verify file extension
	if extension := filepath.Ext(file.Filename); extension != ".csv" {
		result := utils.NewHTTPError(utils.NotCsv)
		utils.Response(ctx, http.StatusBadRequest, result)
		return
	}

	// Update session filename column
	session.FileName = file.Filename
	perr = handler.session.UpdateSession(ctx.Request.Context(), session)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	if err = ctx.SaveUploadedFile(file, handler.filepath+file.Filename); err != nil {
		utils.Response(ctx, http.StatusInternalServerError, utils.NewHTTPError(utils.CouldNotUploadFile))
		return
	}

	result := utils.SuccessPayload(nil, "Successfully uploaded file")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *SessionHandler) DownloadFile(ctx *gin.Context) {
	// Attempt to read from the params
	sessionId, err := uuid.Parse(ctx.Param("sessionId"))
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Guard against cross-tenant download
	organization, _ := middleware.GetOrganizationClaim(ctx)
	session, perr := handler.session.FindById(ctx, sessionId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.SessionNotFound))
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

	file, err := os.Open(handler.filepath + session.FileName)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.FileNotFound))
		return
	}
	defer file.Close()

	buf := &bytes.Buffer{}
	nRead, err := io.Copy(buf, file)
	if err != nil {
		utils.Response(ctx, http.StatusInternalServerError, utils.NewHTTPError(err.Error()))
		return
	}

	ctx.DataFromReader(http.StatusOK, nRead, "text/csv", buf, nil)

	// tempBuffer := make([]byte, 512)
	// ctx.Header("Content-Disposition", "attachment; filename="+session.FileName)
	// ctx.Data(http.StatusOK, "application/octet-stream", tempBuffer)

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
