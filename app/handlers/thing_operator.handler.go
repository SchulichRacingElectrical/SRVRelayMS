package handlers

import (
	"database-ms/app/middleware"
	"database-ms/app/models"
	services "database-ms/app/services"
	utils "database-ms/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ThingOperatorHandler struct {
	service 				services.ThingOperatorServiceInterface
	thingService 		services.ThingServiceInterface
	operatorService services.OperatorServiceInterface
}

func NewThingOperatorAPI(service services.ThingOperatorServiceInterface, thingService services.ThingServiceInterface, operatorService services.OperatorServiceInterface) *ThingOperatorHandler {
	return &ThingOperatorHandler{service: service, thingService: thingService, operatorService: operatorService}
}

func (handler *ThingOperatorHandler) CreateThingOperatorAssociation(ctx *gin.Context) {
	var newThingOperator models.ThingOperator
	ctx.BindJSON(&newThingOperator)
	if handler.IsTenantWithAuth(ctx, "Lead", newThingOperator.ThingId.Hex(),newThingOperator.OperatorId.Hex()) {
		if handler.service.IsAssociationUnique(ctx, &newThingOperator) {
			err := handler.service.Create(ctx.Request.Context(), &newThingOperator)
			if err == nil {
				result := utils.SuccessPayload(nil, "Successfully created association.")
				utils.Response(ctx, http.StatusOK, result)
			} else {
				utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.EntityCreationError))
			}
		} else {
			utils.Response(ctx, http.StatusConflict, utils.NewHTTPError(utils.ThingOperatorNotUnique))
		}
	} else {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
	}
}

func (handler *ThingOperatorHandler) DeleteThingOperator(ctx *gin.Context) {
	if handler.IsTenantWithAuth(ctx, "Lead", ctx.Param("thingId"), ctx.Param("operatorId")) {
		err := handler.service.Delete(ctx, ctx.Param("thingId"), ctx.Param("operatorId"))
		if err == nil {
			result := utils.SuccessPayload(nil, "Thing Operator association successfully deleted.")
			utils.Response(ctx, http.StatusOK, result)
		} else {
			utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		}
	} else {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
	}
}

func (handler *ThingOperatorHandler) IsTenantWithAuth(ctx *gin.Context, authLevel string, thingId string, operatorId string) bool {
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, err := handler.thingService.FindById(ctx, thingId)
	if err != nil {
		utils.Response(ctx, http.StatusNotFound, utils.NewHTTPError(utils.ThingNotFound))
		return false
	}
	operator, err := handler.operatorService.FindById(ctx, operatorId)
	if err != nil {
		utils.Response(ctx, http.StatusNotFound, utils.NewHTTPError(utils.OperatorNotFound))
		return false
	}
	return middleware.IsAuthorizationAtLeast(ctx, authLevel) && thing.OrganizationId == organization.ID && operator.OrganizationId == organization.ID
}