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
		result = utils.SuccessPayload(res, "Successfully created user")
		status = http.StatusOK
	} else {
		fmt.Println(err)
		result = utils.NewHTTPError(utils.EntityCreationError)
		status = http.StatusBadRequest
	}
	utils.Response(c, status, result)
}

// func (handler *UserHandler) FindByUserId(c *gin.Context) {
// 	result := make(map[string]interface{})
// 	user, err := handler.user.FindByUserId(c.Request.Context(), c.Param("userId"))
// 	if err == nil {
// 		result = utils.SuccessPayload(user, "Successfully retrieved user")
// 		utils.Response(c, http.StatusOK, result)
// 	} else {
// 		result = utils.NewHTTPError(utils.UserNotFound)
// 		utils.Response(c, http.StatusBadRequest, result)
// 	}
// }

// func (handler *UserHandler) Update(c *gin.Context) {
// 	var updateUser models.UserUpdate
// 	c.BindJSON(&updateUser)
// 	result := make(map[string]interface{})
// 	err := handler.user.Update(c.Request.Context(), c.Param("userId"), &updateUser)
// 	if err != nil {
// 		result = utils.NewHTTCustomError(utils.BadRequest, err.Error())
// 		utils.Response(c, http.StatusBadRequest, result)
// 		return
// 	}

// 	result = utils.SuccessPayload(nil, "Successfully updated")
// 	utils.Response(c, http.StatusOK, result)
// }

func (handler *UserHandler) Delete(c *gin.Context) {
	result := make(map[string]interface{})
	err := handler.user.Delete(c.Request.Context(), c.Param("userId"))
	if err != nil {
		result = utils.NewHTTCustomError(utils.BadRequest, err.Error())
		utils.Response(c, http.StatusBadRequest, result)
		return
	}

	result = utils.SuccessPayload(nil, "Successfully deleted")
	utils.Response(c, http.StatusOK, result)
}
