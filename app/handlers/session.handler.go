package handlers

import (
	services "database-ms/app/services"

	"github.com/gin-gonic/gin"
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

func (handler *SessionHandler) CreateSession(c *gin.Context) {
	// var newRun models.Run
	// c.BindJSON(&newRun)

	// err := handler.run.CreateRun(c.Request.Context(), &newRun)
	// if err == nil {
	// 	res := &createEntityRes{
	// 		ID: newRun.ID,
	// 	}
	// 	result := utils.SuccessPayload(res, "Successfully created run")
	// 	utils.Response(c, http.StatusOK, result)
	// } else {
	// 	fmt.Println(err)
	// 	result := utils.NewHTTPError(utils.EntityCreationError)
	// 	utils.Response(c, http.StatusBadRequest, result)
	// }
}

func (handler *SessionHandler) GetSessions(c *gin.Context) {
	// 	var run interface{}
	// 	run, err := handler.run.GetRunsByThingId(c.Request.Context(), c.Param("thingId"))

	// 	if err == nil {
	// 		result := utils.SuccessPayload(run, "Successfully retrieved run")
	// 		utils.Response(c, http.StatusOK, result)
	// 	} else {
	// 		result := utils.NewHTTPError(utils.RunsNotFound)
	// 		utils.Response(c, http.StatusBadRequest, result)
	// 	}
	// }

	// func (handler *RunHandler) UpdateRun(c *gin.Context) {
	// 	var updatedRun models.RunUpdate
	// 	c.BindJSON(&updatedRun)

	// 	_, err := handler.run.FindById(c.Request.Context(), updatedRun.ID.Hex())
	// 	if err == nil {
	// 		err := handler.run.UpdateRun(c.Request.Context(), &updatedRun)
	// 		if err == nil {
	// 			result := utils.SuccessPayload(nil, "Successfully updated run.")
	// 			utils.Response(c, http.StatusOK, result)
	// 		} else {
	// 			utils.Response(c, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
	// 		}
	// 	} else {
	// 		utils.Response(c, http.StatusNotFound, utils.NewHTTPError(utils.RunNotFound))
	// 	}
}

func (handler *SessionHandler) UpdateSession(c *gin.Context) {
	// TODO
}

func (handler *SessionHandler) DeleteSession(c *gin.Context) {
	// _, err := handler.run.FindById(c.Request.Context(), c.Param("runId"))
	// if err == nil {
	// 	err := handler.run.DeleteRun(c.Request.Context(), c.Param("runId"))
	// 	if err == nil {
	// 		result := utils.SuccessPayload(nil, "Successfully deleted run.")
	// 		utils.Response(c, http.StatusOK, result)
	// 	} else {
	// 		utils.Response(c, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
	// 	}
	// } else {
	// 	utils.Response(c, http.StatusNotFound, utils.NewHTTPError(utils.RunNotFound))
	// }
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
