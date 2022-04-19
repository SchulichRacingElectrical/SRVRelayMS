package handlers

import (
	middleware "database-ms/app/middleware"
	models "database-ms/app/models"
	services "database-ms/app/services"
	utils "database-ms/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service services.UserServiceInterface
}

func NewUserAPI(userService services.UserServiceInterface) *UserHandler {
	return &UserHandler{service: userService}
}

func (handler *UserHandler) CreateUser(ctx *gin.Context) {
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

func (handler *UserHandler) GetUsers(ctx *gin.Context) {
	organization, err := middleware.GetOrganizationClaim(ctx)
	if err == nil {
		if middleware.IsAuthorizationAtLeast(ctx, "Lead") {
			users, err := handler.service.FindUsersByOrganizationId(ctx.Request.Context(), organization.ID)
			if err == nil {
				result := utils.SuccessPayload(users, "Successfully retrieved users.")
				utils.Response(ctx, http.StatusOK, result)
			} else {
				utils.Response(ctx, http.StatusInternalServerError, "")
			}
		} else {
			utils.Response(ctx, http.StatusUnauthorized, "")
		}
	} else {
		utils.Response(ctx, http.StatusUnauthorized, "")
	}
}

func (handler *UserHandler) Login(ctx *gin.Context) {
	var loggingInUser models.User
	ctx.BindJSON(&loggingInUser)
	result := make(map[string]interface{})

	DBuser, err := handler.service.FindByUserEmail(ctx.Request.Context(), loggingInUser.Email)
	if err != nil {
		result = utils.NewHTTPError(utils.UserNotFound)
		utils.Response(ctx, http.StatusBadRequest, result)
	} else if DBuser.Role == "Pending" {
		result = utils.NewHTTPError(utils.UserNotApproved)
		utils.Response(ctx, http.StatusUnauthorized, result)
	} else {
		if handler.service.CheckPasswordHash(loggingInUser.Password, DBuser.Password) {
			_, err := handler.service.CreateToken(ctx, DBuser)
			if err != nil {
				ctx.JSON(http.StatusUnprocessableEntity, err.Error())
			} else {
				DBuser.Password = ""
				result = utils.SuccessPayload(DBuser, "Successfully signed user in.")
				ctx.JSON(http.StatusOK, result)
			}
		} else {
			result = utils.NewHTTPError(utils.WrongPassword)
			utils.Response(ctx, http.StatusUnauthorized, result)
		}
	}
}

func (handler *UserHandler) UpdateUser(ctx *gin.Context) {
	var updatedUser models.User
	ctx.BindJSON(&updatedUser)
	user, err := middleware.GetUserClaim(ctx)

	if err == nil {
		updatedUser.Role = user.Role // Don't allow users to change their role
		if user.ID == updatedUser.ID {
			err := handler.service.Update(ctx, &updatedUser)
			if err == nil {
				utils.Response(ctx, http.StatusOK, "User updated successfully.")
			} else {
				utils.Response(ctx, http.StatusInternalServerError, "")
			}
		} else {
			utils.Response(ctx, http.StatusUnauthorized, "")	
		}
	} else {
		utils.Response(ctx, http.StatusUnauthorized, "")
	}
}

func (handler *UserHandler) ChangeUserRole(ctx *gin.Context) {
	var updatedUser models.User
	ctx.BindJSON(&updatedUser)
	if middleware.IsAuthorizationAtLeast(ctx, "Lead") {
		err := handler.service.Update(ctx, &updatedUser)
		if err == nil {
			utils.Response(ctx, http.StatusOK, "User promotion successful.")
		} else {
			utils.Response(ctx, http.StatusInternalServerError, "")
		}
	} else {
		utils.Response(ctx, http.StatusUnauthorized, "")
	}
}

func (handler *UserHandler) DeleteUser(ctx *gin.Context) {
	userToDeleteId := ctx.Param("userId")
	user, err := middleware.GetUserClaim(ctx)
	completion := func (ctx *gin.Context, userId string) {
		err := handler.service.Delete(ctx, userToDeleteId)
		if err == nil {
			utils.Response(ctx, http.StatusOK, "User deleted successfully.")
		} else {
			utils.Response(ctx, http.StatusInternalServerError, "")	
		}
	}
	if err == nil {
		if user.ID.String() == userToDeleteId {
			completion(ctx, userToDeleteId)
		} else {
			if middleware.IsAuthorizationAtLeast(ctx, "Admin") {
				completion(ctx, userToDeleteId)
			} else {
				utils.Response(ctx, http.StatusUnauthorized, "")	
			}
		}
	} else {
		utils.Response(ctx, http.StatusUnauthorized, "")
	}
}

func (handler *UserHandler) ChangePassword(ctx *gin.Context) {
	// TODO: Allow the user to change their password
}

func (handler *UserHandler) ForgotPassword(ctx *gin.Context) {
	// TODO: Send an email to the user, somehow flow them to a new password page
	// Likely not worth doing at this time, will take too much time
}