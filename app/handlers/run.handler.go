package handlers

import (
	"database-ms/app/models"
	services "database-ms/app/services"
	utils "database-ms/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RunHandler struct {
	run     services.RunServiceI
	comment services.CommentServiceI
}

func NewRunAPI(runService services.RunServiceI, commentService services.CommentServiceI) *RunHandler {
	return &RunHandler{
		run:     runService,
		comment: commentService,
	}
}

func (handler *RunHandler) CreateRun(c *gin.Context) {
	var newRun models.Run
	c.BindJSON(&newRun)

	err := handler.run.CreateRun(c.Request.Context(), &newRun)
	if err == nil {
		res := &createEntityRes{
			ID: newRun.ID,
		}
		result := utils.SuccessPayload(res, "Successfully created run")
		utils.Response(c, http.StatusOK, result)
	} else {
		fmt.Println(err)
		result := utils.NewHTTPError(utils.EntityCreationError)
		utils.Response(c, http.StatusBadRequest, result)
	}
}

func (handler *RunHandler) GetRuns(c *gin.Context) {
	var run interface{}
	run, err := handler.run.GetRunsByThingId(c.Request.Context(), c.Param("thingId"))

	if err == nil {
		result := utils.SuccessPayload(run, "Successfully retrieved run")
		utils.Response(c, http.StatusOK, result)
	} else {
		result := utils.NewHTTPError(utils.RunsNotFound)
		utils.Response(c, http.StatusBadRequest, result)
	}
}

func (handler *RunHandler) UpdateRun(c *gin.Context) {
	var updatedRun models.RunUpdate
	c.BindJSON(&updatedRun)

	_, err := handler.run.FindById(c.Request.Context(), updatedRun.ID.Hex())
	if err == nil {
		err := handler.run.UpdateRun(c.Request.Context(), &updatedRun)
		if err == nil {
			result := utils.SuccessPayload(nil, "Successfully updated run.")
			utils.Response(c, http.StatusOK, result)
		} else {
			utils.Response(c, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		}
	} else {
		utils.Response(c, http.StatusNotFound, utils.NewHTTPError(utils.RunNotFound))
	}
}

func (handler *RunHandler) DeleteRun(c *gin.Context) {
	_, err := handler.run.FindById(c.Request.Context(), c.Param("runId"))
	if err == nil {
		err := handler.run.DeleteRun(c.Request.Context(), c.Param("runId"))
		if err == nil {
			result := utils.SuccessPayload(nil, "Successfully deleted run.")
			utils.Response(c, http.StatusOK, result)
		} else {
			utils.Response(c, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		}
	} else {
		utils.Response(c, http.StatusNotFound, utils.NewHTTPError(utils.RunNotFound))
	}

}

func (handler *RunHandler) AddComment(c *gin.Context) {
	var comment models.Comment
	c.BindJSON(&comment)

	_, err := handler.run.FindById(c.Request.Context(), c.Param("runId"))
	if err == nil {
		err := handler.comment.AddComment(c.Request.Context(), utils.Run, c.Param("runId"), &comment)
		if err == nil {
			result := utils.SuccessPayload(nil, "Successfully added comment.")
			utils.Response(c, http.StatusOK, result)
		} else {
			utils.Response(c, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		}
	} else {
		utils.Response(c, http.StatusNotFound, utils.NewHTTPError(utils.RunNotFound))
	}
}

func (handler *RunHandler) GetComments(c *gin.Context) {
	comments, err := handler.comment.GetComments(c.Request.Context(), utils.Run, c.Param("runId"))
	if err == nil {
		result := utils.SuccessPayload(comments, "Successfully retrieved comments")
		utils.Response(c, http.StatusOK, result)
	} else {
		utils.Response(c, http.StatusBadRequest, utils.NewHTTPError(utils.CommentsNotFound))
	}
}

func (handler *RunHandler) UpdateCommentContent(c *gin.Context) {
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

func (handler *RunHandler) DeleteComment(c *gin.Context) {
	var requestBody models.Comment
	c.BindJSON(&requestBody)

	if !requestBody.UserID.IsZero() {
		err := handler.comment.DeleteComment(c.Request.Context(), c.Param("commentId"), requestBody.UserID.Hex())
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
