package handlers

import (
	"database-ms/app/middleware"
	"database-ms/app/models"
	services "database-ms/app/services"
	utils "database-ms/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OperatorHandler struct {
	service services.OperatorServiceInterface
}

func NewOperatorAPI(operatorService services.OperatorServiceInterface) *OperatorHandler {
	return &OperatorHandler{service: operatorService}
}

func (handler *OperatorHandler) CreateOperator(ctx *gin.Context) {
	var newOperator models.Operator
	ctx.BindJSON(&newOperator)
	organization, _ := middleware.GetOrganizationClaim(ctx)
	if handler.service.IsOperatorUnique(ctx, &newOperator) {
		if middleware.IsAuthorizationAtLeast(ctx, "Admin") && newOperator.OrganizationId == organization.ID {
			err := handler.service.Create(ctx.Request.Context(), &newOperator)
			if err == nil {
				result := utils.SuccessPayload(newOperator, "Successfully created operator.")
				utils.Response(ctx, http.StatusOK, result)
			} else {
				utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.EntityCreationError))
			}
		} else {
			utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		}
	} else {
		utils.Response(ctx, http.StatusConflict, utils.NewHTTPError(utils.OperatorNotUnique))
	}
}

func (handler *OperatorHandler) GetOperators(ctx *gin.Context) {
	organization, _ := middleware.GetOrganizationClaim(ctx)
	operators, err := handler.service.FindByOrganizationId(ctx.Request.Context(), organization.ID)
	if err == nil {
		result := utils.SuccessPayload(operators, "Successfully retrieved operators.")
		utils.Response(ctx, http.StatusOK, result)
	} else {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingsNotFound))
	}
}

func (handler *OperatorHandler) UpdateOperator(ctx *gin.Context) {
	var operator models.Operator
	ctx.BindJSON(&operator)
	if middleware.IsAuthorizationAtLeast(ctx, "Admin") {
		organization, _ := middleware.GetOrganizationClaim(ctx)
		if organization.ID == operator.OrganizationId {
			err := handler.service.Update(ctx.Request.Context(), &operator)
			if err == nil {
				result := utils.SuccessPayload(nil, "Successfully updated operator.")
				utils.Response(ctx, http.StatusOK, result)
			} else {
				utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
			}
		} else {
			utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))	
		}
	} else {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))	
	}
}

func (handler *OperatorHandler) DeleteOperator(ctx *gin.Context) {
	if middleware.IsAuthorizationAtLeast(ctx, "Admin") {
		organization, _ := middleware.GetOrganizationClaim(ctx)
		operator, err := handler.service.FindById(ctx, ctx.Param("operatorId"))
		if err == nil {
			if organization.ID == operator.OrganizationId {
				err := handler.service.Delete(ctx.Request.Context(), ctx.Param("operatorId"))
				if err == nil {
					result := utils.SuccessPayload(nil, "Successfully deleted operator.")
					utils.Response(ctx, http.StatusOK, result)
				} else {	
					utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
				}
			} else {
				utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))	
			}
		} else {
			utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		}
	} else {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.OperatorNotFound))
	}
}