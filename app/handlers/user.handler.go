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

func (handler *UserHandler) GetUsers(ctx *gin.Context) {
	organization, _ := middleware.GetOrganizationClaim(ctx)
	if middleware.IsAuthorizationAtLeast(ctx, "Lead") {
		users, err := handler.service.FindUsersByOrganizationId(ctx.Request.Context(), organization.ID)
		if err == nil {
			result := utils.SuccessPayload(users, "Successfully retrieved users.")
			utils.Response(ctx, http.StatusOK, result)
		} else {
			utils.Response(ctx, http.StatusNotFound, utils.NewHTTPError(utils.UsersNotFound))
		}
	} else {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
	}
}

func (handler *UserHandler) UpdateUser(ctx *gin.Context) {
	var updatedUser models.User
	ctx.BindJSON(&updatedUser)
	user, err := middleware.GetUserClaim(ctx)
	if err == nil {
		updatedUser.Role = user.Role // Don't allow users to change their role
		if user.ID == updatedUser.ID {
			if handler.service.IsUserUnique(ctx, &updatedUser) {
				err := handler.service.Update(ctx, &updatedUser)
				if err == nil {
					utils.Response(ctx, http.StatusOK, "User updated successfully.")
				} else {
					utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
				}
			} else {
				utils.Response(ctx, http.StatusConflict, utils.NewHTTPError(utils.UserConflict))
			}
		} else {
			utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.InternalError))	
		}
	} else {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.InternalError))
	}
}

func (handler *UserHandler) ChangeUserRole(ctx *gin.Context) {
	var updatedUser models.User
	ctx.BindJSON(&updatedUser)
	organization, _ := middleware.GetOrganizationClaim(ctx)
	if middleware.IsAuthorizationAtLeast(ctx, "Lead") {
		user, err := handler.service.FindByUserId(ctx, updatedUser.ID.Hex())
		if err == nil {
			if user.OrganizationId == organization.ID {
				last, err := handler.service.IsLastAdmin(ctx, user) 
				if err != nil {
					utils.Response(ctx, http.StatusInternalServerError, utils.NewHTTPError(utils.InternalError))
				} else {
					if !last {
						user.Role = updatedUser.Role
						err = handler.service.Update(ctx, user)
						if err == nil {
							utils.Response(ctx, http.StatusOK, "User promotion successful.")
						} else {
							utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
						}
					} else {
						utils.Response(ctx, http.StatusForbidden, utils.NewHTTPError(utils.UserLastAdmin))
					}
				}
			} else {
				utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
			}
		} else {
			utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		}
	} else {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
	}
}

func (handler *UserHandler) DeleteUser(ctx *gin.Context) {
	userToDeleteId := ctx.Param("userId")
	user, err := middleware.GetUserClaim(ctx)
	organization, _ := middleware.GetOrganizationClaim(ctx)
	completion := func (ctx *gin.Context, userId string) {
		err := handler.service.Delete(ctx, userToDeleteId)
		if err == nil {
			utils.Response(ctx, http.StatusOK, "User deleted successfully.")
		} else {
			utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		}
	}
	if err == nil {
		if user.ID.String() == userToDeleteId {
			last, err := handler.service.IsLastAdmin(ctx, user) 
			if err != nil {
				utils.Response(ctx, http.StatusInternalServerError, utils.NewHTTPError(utils.InternalError))
			}
			if !last {
				completion(ctx, userToDeleteId)
			} else {
				utils.Response(ctx, http.StatusForbidden, utils.NewHTTPError(utils.UserLastAdmin))
			}
		} else {
			if middleware.IsAuthorizationAtLeast(ctx, "Admin") {
				user, err := handler.service.FindByUserId(ctx, userToDeleteId)
				if err == nil {
					if user.OrganizationId == organization.ID {
						completion(ctx, userToDeleteId)
					} else {
						utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
					}
				} else {
					utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
				}
			} else {
				utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))	
			}
		}
	} else {
		if middleware.IsAuthorizationAtLeast(ctx, "Admin") { // API Key
			user, err := handler.service.FindByUserId(ctx, userToDeleteId)	
			if err == nil {
				if user.OrganizationId == organization.ID {
					last, err := handler.service.IsLastAdmin(ctx, user) 
					if err != nil {
						utils.Response(ctx, http.StatusInternalServerError, utils.NewHTTPError(utils.InternalError))
					}
					if !last {
						completion(ctx, userToDeleteId)
					} else {
						utils.Response(ctx, http.StatusForbidden, utils.NewHTTPError(utils.UserLastAdmin))
					}
				} else {
					utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
				}
			} else {
				utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
			}
		} else {
			utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		}
	}
}

func (handler *UserHandler) ChangePassword(ctx *gin.Context) {
	// TODO: Allow the user to change their password
}

func (handler *UserHandler) ForgotPassword(ctx *gin.Context) {
	// TODO: Send an email to the user, somehow flow them to a new password page
	// Likely not worth doing at this time, will take too much time
}