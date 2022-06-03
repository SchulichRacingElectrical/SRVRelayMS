package handlers

import (
	"database-ms/app/middleware"
	"database-ms/app/model"
	services "database-ms/app/services"
	utils "database-ms/app/utils"
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
	_, perr := handler.organizationService.FindByOrganizationId(ctx, newUser.OrganizationId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Attempt to get all the users
	users, perr := handler.service.FindUsersByOrganizationId(ctx.Request.Context(), newUser.OrganizationId)
	if perr != nil {
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
	perr = handler.service.Create(ctx.Request.Context(), &newUser)
	if perr != nil {
		if perr.Code == "23505" {
			utils.Response(ctx, http.StatusConflict, utils.NewHTTPError(utils.UserConflict))
		} else {
			utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		}
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
	DBuser, perr := handler.service.FindByUserEmail(ctx.Request.Context(), loggingInUser.Email)
	if perr != nil {
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

func (handler *AuthHandler) Renew(ctx *gin.Context) {
	// Extract token from request
	token, err := middleware.GetToken(ctx)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Blacklist old token
	err = handler.service.BlacklistToken(token)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Extract token from request
	user, err := middleware.GetUserClaim(ctx)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Attempt to create the token
	_, err = handler.service.CreateToken(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	utils.Response(ctx, http.StatusOK, utils.SuccessPayload("", "Successfully renewed token."))
}

func (handler *AuthHandler) Validate(ctx *gin.Context) {
	response := make(map[string]interface{})
	user, _ := middleware.GetUserClaim(ctx)
	if user != nil {
		response["role"] = user.Role
	} else {
		response["role"] = "Admin"
	}
	utils.Response(ctx, http.StatusOK, utils.SuccessPayload(response, "Valid."))
}

func (handler *AuthHandler) SignOut(ctx *gin.Context) {
	// Extract token from request
	token, err := middleware.GetToken(ctx)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Blacklist token
	err = handler.service.BlacklistToken(token)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}
	utils.Response(ctx, http.StatusOK, utils.SuccessPayload("", "Successfully signed user out."))
}
