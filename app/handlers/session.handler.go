package handlers

import (
	"database-ms/app/models"
	"database-ms/app/services"
	utils "database-ms/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	session services.SessionServiceI
	comment services.CommentServiceI
}

func NewSessionAPI(sessionService services.SessionServiceI, commentService services.CommentServiceI) *SessionHandler {
	return &SessionHandler{
		session: sessionService,
		comment: commentService,
	}
}

func (handler *SessionHandler) CreateSession(c *gin.Context) {
	var newSession models.Session
	c.BindJSON(&newSession)

	_, err := handler.session.CreateSession(c.Request.Context(), &newSession)
	if err == nil {
		res := &createEntityRes{
			ID: newSession.ID,
		}
		result := utils.SuccessPayload(res, "Successfully created run")
		utils.Response(c, http.StatusOK, result)
	} else {
		fmt.Println(err)
		result := utils.NewHTTPError(utils.EntityCreationError)
		utils.Response(c, http.StatusBadRequest, result)
	}
}

func (handler *SessionHandler) GetSessions(c *gin.Context) {
	var sessions interface{}
	sessions, err := handler.session.GetSessionsByThingId(c.Request.Context(), c.Param("thingId"))

	if err == nil {
		result := utils.SuccessPayload(sessions, "Successfully retrieved session")
		utils.Response(c, http.StatusOK, result)
	} else {
		result := utils.NewHTTPError(utils.SessionsNotFound)
		utils.Response(c, http.StatusBadRequest, result)
	}
}

func (handler *SessionHandler) UpdateSession(c *gin.Context) {
	var updatedSession models.Session
	c.BindJSON(&updatedSession)

	_, err := handler.session.FindById(c.Request.Context(), updatedSession.ID.Hex())
	if err == nil {
		err := handler.session.UpdateSession(c.Request.Context(), &updatedSession)
		if err == nil {
			result := utils.SuccessPayload(nil, "Successfully updated session.")
			utils.Response(c, http.StatusOK, result)
		} else {
			utils.Response(c, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		}
	} else {
		utils.Response(c, http.StatusNotFound, utils.NewHTTPError(utils.SessionNotFound))
	}
}

func (handler *SessionHandler) DeleteSession(c *gin.Context) {
	_, err := handler.session.FindById(c.Request.Context(), c.Param("sessionId"))
	if err == nil {
		err := handler.session.DeleteSession(c.Request.Context(), c.Param("sessionId"))
		if err == nil {
			result := utils.SuccessPayload(nil, "Successfully deleted session.")
			utils.Response(c, http.StatusOK, result)
		} else {
			utils.Response(c, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		}
	} else {
		utils.Response(c, http.StatusNotFound, utils.NewHTTPError(utils.SessionNotFound))
	}
}

func (handler *SessionHandler) AddComment(c *gin.Context) {
	var comment models.Comment
	c.BindJSON(&comment)

	_, err := handler.session.FindById(c.Request.Context(), c.Param("sessionId"))
	if err == nil {
		err := handler.comment.AddComment(c.Request.Context(), utils.Session, c.Param("sessionId"), &comment)
		if err == nil {
			result := utils.SuccessPayload(nil, "Successfully added comment.")
			utils.Response(c, http.StatusOK, result)
		} else {
			utils.Response(c, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		}
	} else {
		utils.Response(c, http.StatusNotFound, utils.NewHTTPError(utils.SessionNotFound))
	}
}

func (handler *SessionHandler) GetComments(c *gin.Context) {
	comments, err := handler.comment.GetComments(c.Request.Context(), utils.Session, c.Param("sessionId"))
	if err == nil {
		result := utils.SuccessPayload(comments, "Successfully retrieved comments")
		utils.Response(c, http.StatusOK, result)
	} else {
		utils.Response(c, http.StatusBadRequest, utils.NewHTTPError(utils.CommentsNotFound))
	}
}

func (handler *SessionHandler) UpdateCommentContent(c *gin.Context) {
	var updatedComment models.Comment
	c.BindJSON(&updatedComment)

	err := handler.comment.UpdateCommentContent(c.Request.Context(), c.Param("commentId"), &updatedComment)
	if err == nil {
		result := utils.SuccessPayload(nil, "Successfully updated comment")
		utils.Response(c, http.StatusOK, result)
	} else {
		var errMsg string
		switch err.Error() {
		case utils.CommentDoesNotExist, utils.CommentCannotUpdateOtherUserComment:
			errMsg = err.Error()
		default:
			errMsg = utils.BadRequest
		}
		result := utils.NewHTTPError(errMsg)
		utils.Response(c, http.StatusBadRequest, result)
	}
}

func (handler *SessionHandler) DeleteComment(c *gin.Context) {
	var requestBody models.Comment
	c.BindJSON(&requestBody)

	if !requestBody.UserID.IsZero() {
		err := handler.comment.DeleteComment(c.Request.Context(), utils.Session, c.Param("commentId"), requestBody.UserID.Hex())
		if err == nil {
			result := utils.SuccessPayload(nil, "Successfully deleted comment")
			utils.Response(c, http.StatusOK, result)
		} else {
			var errMsg string
			switch err.Error() {
			case utils.CommentDoesNotExist, utils.CommentCannotUpdateOtherUserComment:
				errMsg = err.Error()
			default:
				errMsg = utils.BadRequest
			}
			result := utils.NewHTTPError(errMsg)
			utils.Response(c, http.StatusBadRequest, result)
		}
	} else {
		result := utils.NewHTTPError(utils.UserIdMissing)
		utils.Response(c, http.StatusBadRequest, result)
	}
}
