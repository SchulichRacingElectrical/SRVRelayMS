package handlers

import (
	"database-ms/app/middleware"
	models "database-ms/app/models"
	services "database-ms/app/services"
	utils "database-ms/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ThingHandler struct {
	service services.ThingServiceInterface
}

func NewThingAPI(thingService services.ThingServiceInterface) *ThingHandler {
	return &ThingHandler{service: thingService}
}

func (handler *ThingHandler) CreateThing(ctx *gin.Context) {
	var newThing models.Thing
	ctx.BindJSON(&newThing)
	organization, _ := middleware.GetOrganizationClaim(ctx)	
	if middleware.IsAuthorizationAtLeast(ctx, "Admin") && organization.ID == newThing.OrganizationId {
		err := handler.service.Create(ctx.Request.Context(), &newThing)
		if err == nil {
			result := utils.SuccessPayload(newThing, "Succesfully created thing")
			utils.Response(ctx, http.StatusOK, result)
		} else {
			utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.EntityCreationError))
		}
	} else {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
	}
}

func (handler *ThingHandler) GetThings(ctx *gin.Context) {
	organization, _ := middleware.GetOrganizationClaim(ctx)
	things, err := handler.service.FindByOrganizationId(ctx.Request.Context(), organization.ID)
	if err == nil {
		result := utils.SuccessPayload(things, "Successfully retrieved things.")
		utils.Response(ctx, http.StatusOK, result)
	} else {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingsNotFound))	
	}
}

func (handler *ThingHandler) UpdateThing(ctx *gin.Context) {
	var thing models.Thing
	ctx.BindJSON(&thing)
	if middleware.IsAuthorizationAtLeast(ctx, "Admin") {
		organization, _ := middleware.GetOrganizationClaim(ctx)
		if organization.ID == thing.ID { 
			err := handler.service.Update(ctx.Request.Context(), &thing)
			if err != nil {
				utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
			} else {
				result := utils.SuccessPayload(nil, "Succesfully updated")
				utils.Response(ctx, http.StatusOK, result)
			}
		} else {
			utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))	
		}
	} else {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
	}
}

func (handler *ThingHandler) DeleteThing(ctx *gin.Context) {
	if middleware.IsAuthorizationAtLeast(ctx, "Admin") {
		organization, _ := middleware.GetOrganizationClaim(ctx)
		thing, err := handler.service.FindById(ctx, ctx.Param("thingId"))
		if err == nil {
			if organization.ID == thing.OrganizationId { 
				err := handler.service.Delete(ctx.Request.Context(), ctx.Param("thingId"))
				if err == nil {
					result := utils.SuccessPayload(nil, "Successfully deleted")
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
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
	}
}
