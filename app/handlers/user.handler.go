package handlers

import (
	middleware "database-ms/app/middleware"
	"database-ms/app/model"
	services "database-ms/app/services"
	utils "database-ms/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	service services.UserServiceInterface
}

func NewUserAPI(userService services.UserServiceInterface) *UserHandler {
	return &UserHandler{service: userService}
}

func (handler *UserHandler) GetUsers(ctx *gin.Context) {
	// Guard against non-lead+ requests
	organization, _ := middleware.GetOrganizationClaim(ctx)
	if !middleware.IsAuthorizationAtLeast(ctx, "Lead") {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to get the users
	users, err := handler.service.FindUsersByOrganizationId(ctx.Request.Context(), organization.Id)
	if err != nil {
		utils.Response(ctx, http.StatusNotFound, utils.NewHTTPError(utils.UsersNotFound))
		return
	}

	// Send the response
	result := utils.SuccessPayload(users, "Successfully retrieved users.")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *UserHandler) UpdateUser(ctx *gin.Context) {
	// Attempt to extract the body
	var updatedUser model.User
	err := ctx.BindJSON(&updatedUser)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Attempt to read the user from token
	user, err := middleware.GetUserClaim(ctx)
	if err != nil {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.InternalError))
		return
	}

	// Don't allow users to change their role
	updatedUser.Role = user.Role

	// Guard against a user updated another user's information
	if user.Id != updatedUser.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.InternalError))
		return
	}

	// Attempt to update the user
	perr := handler.service.Update(ctx, &updatedUser)
	if err != nil {
		if perr.Code == "23505" {
			utils.Response(ctx, http.StatusConflict, utils.NewHTTPError(utils.UserConflict))
		} else {
			utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		}
		return
	}

	// Send the response
	utils.Response(ctx, http.StatusOK, "User updated successfully.")
}

func (handler *UserHandler) ChangeUserRole(ctx *gin.Context) {
	// Attempt to extract the body
	var updatedUser model.User
	err := ctx.BindJSON(&updatedUser)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Guard against non-lead+ requestors
	organization, _ := middleware.GetOrganizationClaim(ctx)
	if !middleware.IsAuthorizationAtLeast(ctx, "Lead") {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to find the existing user
	user, perr := handler.service.FindByUserId(ctx, updatedUser.Id)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Guard against cross-tenant roles changes
	if user.OrganizationId != organization.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Guard against removing the last admin
	last, err := handler.service.IsLastAdmin(ctx, user)
	if err != nil {
		utils.Response(ctx, http.StatusInternalServerError, utils.NewHTTPError(utils.InternalError))
		return
	}
	if last {
		// TODO: This should say "Cannot demote last admin"
		utils.Response(ctx, http.StatusForbidden, utils.NewHTTPError(utils.UserLastAdmin))
		return
	}

	// Attempt to update the user's role
	user.Role = updatedUser.Role
	perr = handler.service.Update(ctx, user)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Send the response
	utils.Response(ctx, http.StatusOK, "User promotion successful.")
}

func (handler *UserHandler) DeleteUser(ctx *gin.Context) {
	// Completion on deletion
	completion := func(ctx *gin.Context, userID uuid.UUID) {
		// Attempt to delete the user
		err := handler.service.Delete(ctx, userID)
		if err != nil {
			utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
			return
		}

		// Send the response
		utils.Response(ctx, http.StatusOK, "User deleted successfully.")
	}

	// Attempt to parse the query param
	userIdToDelete, err := uuid.Parse(ctx.Param("userId"))
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Handle deletion from user request
	user, err := middleware.GetUserClaim(ctx)
	organization, _ := middleware.GetOrganizationClaim(ctx)
	if err == nil {
		// Only allow the user to delete if the request comes from themselves
		if user.Id == userIdToDelete {
			// Guard against deletion of last admin user
			last, err := handler.service.IsLastAdmin(ctx, user)
			if err != nil {
				utils.Response(ctx, http.StatusInternalServerError, utils.NewHTTPError(utils.InternalError))
				return
			}
			if last {
				utils.Response(ctx, http.StatusForbidden, utils.NewHTTPError(utils.UserLastAdmin))
				return
			}

			// Send the response
			completion(ctx, userIdToDelete)
			return
		}

		// Guard against deletion from non-admin requestors
		if !middleware.IsAuthorizationAtLeast(ctx, "Admin") {
			utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
			return
		}

		// Attempt to find the user
		user, perr := handler.service.FindByUserId(ctx, userIdToDelete)
		if perr != nil {
			utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
			return
		}

		// Guard against cross-tenant deletions
		if user.OrganizationId != organization.Id {
			utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
			return
		}

		// Send the response
		completion(ctx, userIdToDelete)
		return
	}

	// Handle deletion from API Key, Guard against non-admin requests
	if !middleware.IsAuthorizationAtLeast(ctx, "Admin") {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to find the user
	user, perr := handler.service.FindByUserId(ctx, userIdToDelete)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Guard against cross-tenant deletion
	if user.OrganizationId != organization.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Guard against deleting the last admin
	last, err := handler.service.IsLastAdmin(ctx, user)
	if err != nil {
		utils.Response(ctx, http.StatusInternalServerError, utils.NewHTTPError(utils.InternalError))
		return
	}
	if last {
		utils.Response(ctx, http.StatusForbidden, utils.NewHTTPError(utils.UserLastAdmin))
		return
	}

	// Send the response
	completion(ctx, userIdToDelete)
}

func (handler *UserHandler) ChangePassword(ctx *gin.Context) {
	// TODO: Allow the user to change their password
}

func (handler *UserHandler) ForgotPassword(ctx *gin.Context) {
	// TODO: Send an email to the user, somehow flow them to a new password page
	// Likely not worth doing at this time, will take too much time
}
