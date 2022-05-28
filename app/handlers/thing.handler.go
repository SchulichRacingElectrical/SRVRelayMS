package handlers

import (
	"database-ms/app/middleware"
	"database-ms/app/model"
	services "database-ms/app/services"
	utils "database-ms/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ThingHandler struct {
	service services.ThingServiceInterface
}

func NewThingAPI(thingService services.ThingServiceInterface) *ThingHandler {
	return &ThingHandler{service: thingService}
}

func (handler *ThingHandler) CreateThing(ctx *gin.Context) {
	// Attempt to extract the body
	var newThing model.Thing
	err := ctx.BindJSON(&newThing)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Guard against non-unique things
	organization, _ := middleware.GetOrganizationClaim(ctx)
	newThing.OrganizationId = organization.Id
	if !handler.service.IsThingUnique(ctx, &newThing) {
		utils.Response(ctx, http.StatusConflict, utils.NewHTTPError(utils.ThingNotUnique))
		return
	}

	// Guard against non-admin requests
	if !middleware.IsAuthorizationAtLeast(ctx, "Admin") {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to create the thing
	err = handler.service.Create(ctx.Request.Context(), &newThing)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.EntityCreationError))
		return
	}

	// Send the response
	result := utils.SuccessPayload(newThing, "Succesfully created thing.")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *ThingHandler) GetThings(ctx *gin.Context) {
	// Attempt to fetch the things
	organization, _ := middleware.GetOrganizationClaim(ctx)
	things, err := handler.service.FindByOrganizationId(ctx.Request.Context(), organization.Id)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingsNotFound))
		return
	}

	// Send the response
	result := utils.SuccessPayload(things, "Successfully retrieved things.")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *ThingHandler) UpdateThing(ctx *gin.Context) {
	// Attempt to extract the body
	var updatedThing model.Thing
	err := ctx.BindJSON(&updatedThing)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Guard against non-admin users
	if !middleware.IsAuthorizationAtLeast(ctx, "Admin") {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to find the thing we are updated
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, err := handler.service.FindById(ctx, updatedThing.Id)
	if err != nil {
		utils.Response(ctx, http.StatusNotFound, utils.NewHTTPError(utils.ThingNotFound))
		return
	}

	// Guard against cross-tenant writes
	if organization.Id != thing.OrganizationId {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Guard against non-unique things
	updatedThing.OrganizationId = thing.OrganizationId
	if !handler.service.IsThingUnique(ctx, &updatedThing) {
		utils.Response(ctx, http.StatusConflict, utils.NewHTTPError(utils.ThingNotUnique))
		return
	}

	// Attempt to update the thing
	err = handler.service.Update(ctx.Request.Context(), &updatedThing)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Send the response
	result := utils.SuccessPayload(nil, "Succesfully updated thing.")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *ThingHandler) DeleteThing(ctx *gin.Context) {
	// Attempt to read from the params
	thingIDToDelete, err := uuid.FromBytes([]byte(ctx.Param("thingId")))
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Guard against non-admin requests
	if !middleware.IsAuthorizationAtLeast(ctx, "Admin") {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.ThingNotFound))
		return
	}

	// Attempt to find the thing
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, err := handler.service.FindById(ctx, thingIDToDelete)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Guard against cross-tenant deletion
	if organization.Id != thing.OrganizationId {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to delete the thing
	err = handler.service.Delete(ctx.Request.Context(), thingIDToDelete)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Send the response
	result := utils.SuccessPayload(nil, "Successfully deleted thing.")
	utils.Response(ctx, http.StatusOK, result)
}
