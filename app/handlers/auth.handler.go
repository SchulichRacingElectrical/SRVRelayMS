package handlers

import (
	models "database-ms/app/models"
	services "database-ms/app/services"
	utils "database-ms/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service services.UserServiceInterface
	// Add a specific auth service
}

func NewAuthAPI(userService services.UserServiceInterface) *UserHandler {
	return &UserHandler{service: userService}
}

func (handler *UserHandler) SignUp(ctx *gin.Context) {
	var newUser models.User
	ctx.BindJSON(&newUser)
	if !handler.service.IsUserUnique(ctx.Request.Context(), &newUser) {
		result := utils.NewHTTPError(utils.UserAlreadyExists)
		utils.Response(ctx, http.StatusConflict, result)
	} else {
		users, err := handler.service.FindUsersByOrganizationId(ctx.Request.Context(), newUser.OrganizationId)
		if err == nil { 
			// If there are no users in the organization, set the first user as an admin
			if len(users) == 0 {
				newUser.Role = "Admin"
			} else {
				newUser.Role = "Pending"
			}
			newUser.Password = handler.service.HashPassword(newUser.Password)
			_, err := handler.service.Create(ctx.Request.Context(), &newUser)
			if err == nil {
				newUser.Password = ""
				result := utils.SuccessPayload(newUser, "Successfully created user.")
				utils.Response(ctx, http.StatusOK, result)
			} else {
				result := utils.NewHTTPError(utils.EntityCreationError)
				utils.Response(ctx, http.StatusBadRequest, result)
			}
		} else {
			result := utils.SuccessPayload("", "Invalid Organization.")
			utils.Response(ctx, http.StatusBadRequest, result)
		}
	}
}

func (handler *UserHandler) Login(ctx *gin.Context) {
	var loggingInUser models.User
	ctx.BindJSON(&loggingInUser)
	DBuser, err := handler.service.FindByUserEmail(ctx.Request.Context(), loggingInUser.Email)
	if err != nil {
		result := utils.NewHTTPError(utils.UserNotFound)
		utils.Response(ctx, http.StatusBadRequest, result)
	} else if DBuser.Role == "Pending" {
		result := utils.NewHTTPError(utils.UserNotApproved)
		utils.Response(ctx, http.StatusUnauthorized, result)
	} else {
		if handler.service.CheckPasswordHash(loggingInUser.Password, DBuser.Password) {
			_, err := handler.service.CreateToken(ctx, DBuser)
			if err != nil {
				ctx.JSON(http.StatusUnprocessableEntity, err.Error())
			} else {
				DBuser.Password = ""
				result := utils.SuccessPayload(DBuser, "Successfully signed user in.")
				ctx.JSON(http.StatusOK, result)
			}
		} else {
			result := utils.NewHTTPError(utils.WrongPassword)
			utils.Response(ctx, http.StatusUnauthorized, result)
		}
	}
}

func (handler *UserHandler) Validate(ctx *gin.Context) {
	utils.Response(ctx, http.StatusOK, "Valid.")
}

func (handler *UserHandler) SignOut(ctx *gin.Context) {
	// TODO: Blacklist tokens
	// TODO: Delete blacklisted tokens in the database after they expire
	// TODO: In auth middleware, check if the token is blacklisted
}