package handlers

import (
	"database-ms/app/models"
	"database-ms/app/services/run"
	"database-ms/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RunHandler struct {
	run run.RunServiceI
}

func NewRunAPI(runService run.RunServiceI) *RunHandler {
	return &RunHandler{
		run: runService,
	}
}

func (handler *RunHandler) Create(c *gin.Context) {
	var newRun models.Run
	c.BindJSON(&newRun)
	var result map[string]interface{}

	err := handler.run.Create(c.Request.Context(), &newRun)
	var status int
	if err == nil {
		res := &createEntityRes{
			ID: newRun.ID,
		}
		result = utils.SuccessPayload(res, "Successfully created sensor")
		status = http.StatusOK
	} else {
		fmt.Println(err)
		result = utils.NewHTTPError(utils.EntityCreationError)
		status = http.StatusBadRequest
	}
	utils.Response(c, status, result)
}

func (handler *RunHandler) GetRun(c *gin.Context) {
	var result map[string]interface{}

	var run interface{}
	run, err := handler.run.GetRun(c.Request.Context(), c.Param("runId"))

	if err == nil {
		result = utils.SuccessPayload(run, "Successfully retrieved run")
		utils.Response(c, http.StatusOK, result)
	} else {
		result = utils.NewHTTPError(utils.RunNotFound)
		utils.Response(c, http.StatusBadRequest, result)
	}
}

func (handler *RunHandler) GetComments(c *gin.Context) {
	var result map[string]interface{}

	comments, err := handler.run.GetComments(c.Request.Context(), c.Param("runId"))
	if err == nil {
		result = utils.SuccessPayload(comments, "Successfully retrieved comments")
		utils.Response(c, http.StatusOK, result)
	} else {
		result = utils.NewHTTPError(utils.CommentsNotFound)
		utils.Response(c, http.StatusBadRequest, result)
	}

}
