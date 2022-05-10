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
}

func NewSessionAPI(sessionService services.SessionServiceI) *SessionHandler {
	return &SessionHandler{
		session: sessionService,
	}
}

func (handler *SessionHandler) CreateSession(c *gin.Context) {
	var newSession models.Session
	c.BindJSON(&newSession)

	err := handler.session.CreateSession(c.Request.Context(), &newSession)
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
