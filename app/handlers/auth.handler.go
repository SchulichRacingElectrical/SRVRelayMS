package handlers

import (
	"database-ms/app/model"
	services "database-ms/app/services"
	utils "database-ms/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service             services.UserServiceInterface
	organizationService services.OrganizationServiceInterface
}

func NewAuthAPI(userService services.UserServiceInterface, organizationService services.OrganizationServiceInterface) *AuthHandler {
	return &AuthHandler{service: userService, organizationService: organizationService}
}

func (handler *AuthHandler) SignUp(ctx *gin.Context) {
	// Attempt to extract the body
	var newUser model.User
	err := ctx.BindJSON(&newUser)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Ensure the organization exists
	_, err = handler.organizationService.FindByOrganizationId(ctx, newUser.OrganizationId)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Ensure the user is unique
	if !handler.service.IsUserUnique(ctx.Request.Context(), &newUser) {
		result := utils.NewHTTPError(utils.UserConflict)
		utils.Response(ctx, http.StatusConflict, result)
		return
	}

	// Attempt to get all the users
	users, err := handler.service.FindUsersByOrganizationId(ctx.Request.Context(), newUser.OrganizationId)
	if err != nil {
		result := utils.SuccessPayload("", "Invalid Organization.")
		utils.Response(ctx, http.StatusBadRequest, result)
		return
	}

	// If there are no users in the organization, set the first user as an admin
	if len(users) == 0 {
		newUser.Role = "Admin"
	} else {
		newUser.Role = "Pending"
	}
	newUser.Password = handler.service.HashPassword(newUser.Password)

	// Attempt to create the user
	_, err = handler.service.Create(ctx.Request.Context(), &newUser)
	if err != nil {
		result := utils.NewHTTPError(utils.EntityCreationError)
		utils.Response(ctx, http.StatusBadRequest, result)
		return
	}

	// Send the response
	newUser.Password = ""
	result := utils.SuccessPayload(newUser, "Successfully created user.")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *AuthHandler) Login(ctx *gin.Context) {
	// Attempt to extract the body
	var loggingInUser model.User
	err := ctx.BindJSON(&loggingInUser)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Check if the user exists
	DBuser, err := handler.service.FindByUserEmail(ctx.Request.Context(), loggingInUser.Email)
	if err != nil {
		result := utils.NewHTTPError(utils.UserNotFound)
		utils.Response(ctx, http.StatusBadRequest, result)
		return
	}

	// Guard against pending users
	if DBuser.Role == "Pending" {
		result := utils.NewHTTPError(utils.UserNotApproved)
		utils.Response(ctx, http.StatusUnauthorized, result)
		return
	}

	// Check if the password matches
	if !handler.service.CheckPasswordHash(loggingInUser.Password, DBuser.Password) {
		result := utils.NewHTTPError(utils.WrongPassword)
		utils.Response(ctx, http.StatusUnauthorized, result)
		return
	}

	// Attempt to create the token
	_, err = handler.service.CreateToken(ctx, DBuser)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Send the response
	DBuser.Password = ""
	result := utils.SuccessPayload(DBuser, "Successfully signed user in.")
	ctx.JSON(http.StatusOK, result)
}

func (handler *AuthHandler) Validate(ctx *gin.Context) {
	// TODO: Send back the authorization level
	utils.Response(ctx, http.StatusOK, "Valid.")
}

func (handler *AuthHandler) SignOut(ctx *gin.Context) {
	// TODO: Blacklist tokens
	// TODO: Delete blacklisted tokens in the database after they expire
	// TODO: In auth middleware, check if the token is blacklisted
}
