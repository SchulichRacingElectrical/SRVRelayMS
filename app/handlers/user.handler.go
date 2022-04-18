package handlers

import (
	"database-ms/app/models"
	services "database-ms/app/services"
	"database-ms/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service services.UserServiceInterface
}

func NewUserAPI(userService services.UserServiceInterface) *UserHandler {
	return &UserHandler{service: userService}
}

func (handler *UserHandler) Create(ctx *gin.Context) {
	var newUser models.User
	ctx.BindJSON(&newUser)
	result := make(map[string]interface{})

	if !handler.service.IsUserUnique(ctx.Request.Context(), &newUser) {
		result = utils.NewHTTPError(utils.UserAlreadyExists)
		utils.Response(ctx, http.StatusConflict, result)
		return
	}

	// If there are no users, set the first user as an admin
	users, err := handler.service.FindUsersByOrganizationId(ctx.Request.Context(), newUser.OrganizationId)
	if err != nil { 
		result = utils.SuccessPayload("", "Invalid Organization.")
		utils.Response(ctx, http.StatusBadRequest, result)
	} else {
		if len(users) == 0 {
			newUser.Role = "Admin"
		} else {
			newUser.Role = "Pending"
		}
	}

	newUser.Password = handler.service.HashPassword(newUser.Password)
	res, err := handler.service.Create(ctx.Request.Context(), &newUser)
	if err == nil {
		result = utils.SuccessPayload(res, "Successfully created user")
		utils.Response(ctx, http.StatusOK, result)
	} else {
		result = utils.NewHTTPError(utils.EntityCreationError)
		utils.Response(ctx, http.StatusBadRequest, result)
	}
}

func (handler *UserHandler) GetUsers(ctx *gin.Context) {
	// Create auth functions to check permissions
	userInterface, _ := ctx.Get("user")
	user := userInterface.(*models.User)
	permitted := user.Role == "Admin" || user.Role == "Lead"
	if permitted {
		users, err := handler.service.FindUsersByOrganizationId(ctx.Request.Context(), user.OrganizationId)
		if err != nil {
			utils.Response(ctx, http.StatusInternalServerError, "")
		} else {
			result := utils.SuccessPayload(users, "Successfully retrieved users.")
			utils.Response(ctx, http.StatusOK, result)
		}
	} else {
		utils.Response(ctx, http.StatusBadRequest, "")
	}
}

func (handler *UserHandler) Login(c *gin.Context) {
	var loggingInUser models.User
	c.BindJSON(&loggingInUser)
	result := make(map[string]interface{})

	DBuser, err := handler.service.FindByUserEmail(c.Request.Context(), loggingInUser.Email)
	if err != nil {
		result = utils.NewHTTPError(utils.UserNotFound)
		utils.Response(c, http.StatusBadRequest, result)
		return
	} else if DBuser.Role == "Pending" {
		result = utils.NewHTTPError(utils.UserNotApproved)
		utils.Response(c, http.StatusUnauthorized, result)
		return
	}

	// Check password match
	if handler.service.CheckPasswordHash(loggingInUser.Password, DBuser.Password) {
		_, err := handler.service.CreateToken(c, DBuser)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, err.Error())
			return
		}
		DBuser.Password = ""
		result = utils.SuccessPayload(DBuser, "Successfully signed user in.")
		c.JSON(http.StatusOK, result)
	} else {
		result = utils.NewHTTPError(utils.WrongPassword)
		utils.Response(c, http.StatusUnauthorized, result)
	}
}

