package handlers

import (
	"database-ms/app/models"
	userSrv "database-ms/app/services/user"
	"database-ms/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	user userSrv.UserServiceInterface
}

func NewUserAPI(userService userSrv.UserServiceInterface) *UserHandler {
	return &UserHandler{
		user: userService,
	}
}

func (handler *UserHandler) Create(c *gin.Context) {
	var newUser models.User
	c.BindJSON(&newUser)
	result := make(map[string]interface{})

	err := handler.user.Create(c.Request.Context(), &newUser)
	var status int
	if err == nil {
		res := &createUserRes{
			ID: newUser.ID,
		}
		result = utils.SuccessPayload(res, "Succesfully created user")
		status = http.StatusOK
	} else {
		fmt.Println(err)
		result = utils.NewHTTPError(utils.EntityCreationError)
		status = http.StatusBadRequest
	}
	utils.Response(c, status, result)
}

func (handler *UserHandler) GetUser(c *gin.Context) {
	// result := make(map[string]interface{})
	// users, err := handler.user.GetUsers(c.Request.Context(), c.Param(""))
	// if err == nil {
	// 	result = utils.SuccessPayload(users, "Succesfully retrieved users")
	// 	utils.Response(c, http.StatusOK, result)
	// } else {
	// 	result = utils.NewHTTPError(utils.SensorsNotFound)
	// 	utils.Response(c, http.StatusBadRequest, result)
	// }
}

func (handler *UserHandler) GetUsers(c *gin.Context) {
	result := make(map[string]interface{})
	users, err := handler.user.GetUsers(c.Request.Context(), c.Param(""))
	if err == nil {
		result = utils.SuccessPayload(users, "Succesfully retrieved users")
		utils.Response(c, http.StatusOK, result)
	} else {
		result = utils.NewHTTPError(utils.SensorsNotFound)
		utils.Response(c, http.StatusBadRequest, result)
	}
}
